package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	flag.Parse()

	log.Printf("Starting stress test with %d workers for %v", *workers, *duration)
	log.Printf("Target: %s", *baseURL)

	if err := setupInitialData(); err != nil {
		log.Fatalf("Failed to setup initial data: %v", err)
	}

	startTime := time.Now()
	stopChan := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			runWorker(workerID, stopChan)
		}(i)
	}

	go printStats(startTime, stopChan)

	time.Sleep(*duration)
	close(stopChan)
	wg.Wait()

	elapsed := time.Since(startTime)
	total := requestsTotal.Load()
	ok := requestsOK.Load()
	fail := requestsFail.Load()

	log.Printf("\n=== Final Results ===")
	log.Printf("Duration: %v", elapsed)
	log.Printf("Total requests: %d", total)
	log.Printf("Successful: %d (%.2f%%)", ok, float64(ok)/float64(total)*100)
	log.Printf("Failed: %d (%.2f%%)", fail, float64(fail)/float64(total)*100)
	log.Printf("RPS: %.2f", float64(total)/elapsed.Seconds())
}

func setupInitialData() error {
	log.Println("Setting up initial data...")

	for _, teamName := range teamNames {
		members := make([]TeamMember, 0)
		for i := 0; i < 40; i++ {
			userID := fmt.Sprintf("%s-u%d", teamName, i)
			userIDs = append(userIDs, userID)
			members = append(members, TeamMember{
				UserID:   userID,
				Username: fmt.Sprintf("User %s %d", teamName, i),
				IsActive: true,
			})
		}

		team := Team{
			TeamName: teamName,
			Members:  members,
		}

		if err := createTeam(team); err != nil {
			return fmt.Errorf("failed to create team %s: %w", teamName, err)
		}
	}

	log.Printf("Created %d teams with %d users", len(teamNames), len(userIDs))
	return nil
}

func createTeam(team Team) error {
	data, err := json.Marshal(team)
	if err != nil {
		return err
	}

	resp, err := http.Post(*baseURL+"/team/add", "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	return nil
}

func runWorker(workerID int, stopChan <-chan struct{}) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))

	for {
		select {
		case <-stopChan:
			return
		default:
			action := rng.Intn(100)

			switch {
			case action < 40:
				createPR(rng)
			case action < 50:
				mergePR(rng)
			case action < 60:
				reassignPR(rng)
			case action < 70:
				getUserReviews(rng)
			case action < 80:
				getTeam(rng)
			case action < 85:
				setUserIsActive(rng)
			case action < 90:
				safeReassignPR(rng)
			case action < 95:
				getStats()
			default:
				deactivateTeam(rng)
			}
		}
	}
}

func createPR(rng *rand.Rand) {
	if len(userIDs) == 0 {
		return
	}

	prID := fmt.Sprintf("pr-%d-%d", time.Now().UnixNano(), rng.Intn(1000000))
	authorID := userIDs[rng.Intn(len(userIDs))]

	req := CreatePRRequest{
		PullRequestID:   prID,
		PullRequestName: fmt.Sprintf("Feature %s", prID),
		AuthorID:        authorID,
	}

	data, _ := json.Marshal(req)
	resp, err := http.Post(*baseURL+"/pullRequest/create", "application/json", bytes.NewReader(data))

	requestsTotal.Add(1)
	if err != nil || resp.StatusCode != http.StatusCreated {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()

	requestsOK.Add(1)

	prIDsMutex.Lock()
	prIDs = append(prIDs, prID)
	if len(prIDs) > 1000 {
		prIDs = prIDs[len(prIDs)-1000:]
	}
	prIDsMutex.Unlock()
}

func mergePR(rng *rand.Rand) {
	prIDsMutex.RLock()
	if len(prIDs) == 0 {
		prIDsMutex.RUnlock()
		return
	}
	prID := prIDs[rng.Intn(len(prIDs))]
	prIDsMutex.RUnlock()

	req := MergePRRequest{PullRequestID: prID}
	data, _ := json.Marshal(req)

	resp, err := http.Post(*baseURL+"/pullRequest/merge", "application/json", bytes.NewReader(data))

	requestsTotal.Add(1)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound) {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func reassignPR(rng *rand.Rand) {
	prIDsMutex.RLock()
	if len(prIDs) == 0 || len(userIDs) == 0 {
		prIDsMutex.RUnlock()
		return
	}
	prID := prIDs[rng.Intn(len(prIDs))]
	prIDsMutex.RUnlock()

	oldUserID := userIDs[rng.Intn(len(userIDs))]
	req := ReassignPRRequest{
		PullRequestID: prID,
		OldUserID:     oldUserID,
	}

	data, _ := json.Marshal(req)
	resp, err := http.Post(*baseURL+"/pullRequest/reassign", "application/json", bytes.NewReader(data))

	requestsTotal.Add(1)
	if err != nil {
		requestsFail.Add(1)
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusConflict {
		requestsOK.Add(1)
	} else {
		requestsFail.Add(1)
	}
}

func getUserReviews(rng *rand.Rand) {
	if len(userIDs) == 0 {
		return
	}

	userID := userIDs[rng.Intn(len(userIDs))]
	url := fmt.Sprintf("%s/users/getReview?user_id=%s", *baseURL, userID)

	resp, err := http.Get(url)

	requestsTotal.Add(1)
	if err != nil || resp.StatusCode != http.StatusOK {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func getTeam(rng *rand.Rand) {
	if len(teamNames) == 0 {
		return
	}

	teamName := teamNames[rng.Intn(len(teamNames))]
	url := fmt.Sprintf("%s/team/get?team_name=%s", *baseURL, teamName)

	resp, err := http.Get(url)

	requestsTotal.Add(1)
	if err != nil || resp.StatusCode != http.StatusOK {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func setUserIsActive(rng *rand.Rand) {
	if len(userIDs) == 0 {
		return
	}

	userID := userIDs[rng.Intn(len(userIDs))]
	req := SetIsActiveRequest{
		UserID:   userID,
		IsActive: rng.Intn(2) == 0,
	}

	data, _ := json.Marshal(req)
	resp, err := http.Post(*baseURL+"/users/setIsActive", "application/json", bytes.NewReader(data))

	requestsTotal.Add(1)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound) {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func safeReassignPR(rng *rand.Rand) {
	prIDsMutex.RLock()
	if len(prIDs) == 0 {
		prIDsMutex.RUnlock()
		return
	}
	prID := prIDs[rng.Intn(len(prIDs))]
	prIDsMutex.RUnlock()

	req := SafeReassignRequest{PullRequestID: prID}
	data, _ := json.Marshal(req)

	resp, err := http.Post(*baseURL+"/pullRequest/safeReassign", "application/json", bytes.NewReader(data))

	requestsTotal.Add(1)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound) {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func getStats() {
	resp, err := http.Get(*baseURL + "/stats")

	requestsTotal.Add(1)
	if err != nil || resp.StatusCode != http.StatusOK {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func deactivateTeam(rng *rand.Rand) {
	if len(teamNames) == 0 {
		return
	}

	teamName := teamNames[rng.Intn(len(teamNames))]
	url := fmt.Sprintf("%s/teams/%s/deactivate", *baseURL, teamName)

	resp, err := http.Post(url, "application/json", nil)

	requestsTotal.Add(1)
	if err != nil || (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound) {
		requestsFail.Add(1)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return
	}
	defer func() { _ = resp.Body.Close() }()
	requestsOK.Add(1)
}

func printStats(startTime time.Time, stopChan <-chan struct{}) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	lastTotal := uint64(0)
	lastTime := startTime

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			now := time.Now()
			total := requestsTotal.Load()
			ok := requestsOK.Load()
			fail := requestsFail.Load()

			elapsed := now.Sub(startTime)
			intervalDuration := now.Sub(lastTime)
			intervalRequests := total - lastTotal

			currentRPS := float64(intervalRequests) / intervalDuration.Seconds()
			avgRPS := float64(total) / elapsed.Seconds()

			log.Printf("[%v] Total: %d | OK: %d | Fail: %d | Current RPS: %.2f | Avg RPS: %.2f",
				elapsed.Round(time.Second), total, ok, fail, currentRPS, avgRPS)

			lastTotal = total
			lastTime = now
		}
	}
}

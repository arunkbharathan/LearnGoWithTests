package poker_test

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	poker "github.com/arunkbharathan/learnWithTests/time"
)

var dummyBlindAlerter = &SpyBlindAlerter{}
var dummyPlayerStore = &poker.StubPlayerStore{}
var dummyStdIn = &bytes.Buffer{}
var dummyStdOut = &bytes.Buffer{}

// for input
const PlayerPrompt = "Please enter the number of players: "

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

var dummySpyAlerter = &SpyBlindAlerter{}

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(5)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 10 * time.Minute, Amount: 200},
			{At: 20 * time.Minute, Amount: 300},
			{At: 30 * time.Minute, Amount: 400},
			{At: 40 * time.Minute, Amount: 500},
			{At: 50 * time.Minute, Amount: 600},
			{At: 60 * time.Minute, Amount: 800},
			{At: 70 * time.Minute, Amount: 1000},
			{At: 80 * time.Minute, Amount: 2000},
			{At: 90 * time.Minute, Amount: 4000},
			{At: 100 * time.Minute, Amount: 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &poker.SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		game.Start(7)

		cases := []poker.ScheduledAlert{
			{At: 0 * time.Second, Amount: 100},
			{At: 12 * time.Minute, Amount: 200},
			{At: 24 * time.Minute, Amount: 300},
			{At: 36 * time.Minute, Amount: 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

}

func TestGame_Finish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	game := poker.NewGame(dummyBlindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	poker.AssertPlayerWin(t, store, winner)
}
func TestCLI(t *testing.T) {
	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		cli := poker.NewCLI(in, stdout, game)

		cli.PlayPoker()

		got := stdout.String()
		want := poker.PlayerPrompt

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {

				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
	// t.Run("it schedules printing of blind values", func(t *testing.T) {
	// 	in := strings.NewReader("Chris wins\n")
	// 	playerStore := &poker.StubPlayerStore{}
	// 	blindAlerter := &SpyBlindAlerter{}
	// 	game := poker.NewGame(blindAlerter, playerStore)

	// 	cli := poker.NewCLI(in, dummyStdOut, game)
	// 	cli.PlayPoker()

	// 	cases := []scheduledAlert{
	// 		{0 * time.Second, 100},
	// 		{10 * time.Minute, 200},
	// 		{20 * time.Minute, 300},
	// 		{30 * time.Minute, 400},
	// 		{40 * time.Minute, 500},
	// 		{50 * time.Minute, 600},
	// 		{60 * time.Minute, 800},
	// 		{70 * time.Minute, 1000},
	// 		{80 * time.Minute, 2000},
	// 		{90 * time.Minute, 4000},
	// 		{100 * time.Minute, 8000},
	// 	}

	// 	for i, want := range cases {
	// 		t.Run(fmt.Sprint(want), func(t *testing.T) {

	// 			if len(blindAlerter.alerts) <= i {
	// 				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
	// 			}

	// 			got := blindAlerter.alerts[i]
	// 			assertScheduledAlert(t, got, want)
	// 		})
	// 	}
	// })

	// t.Run("record chris win from user input", func(t *testing.T) {
	// 	in := strings.NewReader("Chris wins\n")
	// 	playerStore := &poker.StubPlayerStore{}
	// 	game := poker.NewGame(dummySpyAlerter, playerStore)

	// 	cli := poker.NewCLI(in, dummyStdOut, game)
	// 	cli.PlayPoker()

	// 	poker.AssertPlayerWin(t, playerStore, "Chris")
	// })

	// t.Run("record cleo win from user input", func(t *testing.T) {
	// 	in := strings.NewReader("Cleo wins\n")
	// 	playerStore := &poker.StubPlayerStore{}
	// 	game := poker.NewGame(dummySpyAlerter, playerStore)

	// 	cli := poker.NewCLI(in, dummyStdOut, game)
	// 	cli.PlayPoker()

	// 	poker.AssertPlayerWin(t, playerStore, "Cleo")
	// })

	t.Run("do not read beyond the first newline", func(t *testing.T) {
		in := failOnEndReader{
			t,
			strings.NewReader("Chris wins\n hello there"),
		}

		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummySpyAlerter, playerStore)

		cli := poker.NewCLI(in, dummyStdOut, game)
		cli.PlayPoker()
	})

}

type failOnEndReader struct {
	t   *testing.T
	rdr io.Reader
}

func (m failOnEndReader) Read(p []byte) (n int, err error) {

	n, err = m.rdr.Read(p)

	if n == 0 || err == io.EOF {
		m.t.Fatal("Read to the end when you shouldn't have")
	}

	return n, err
}

func assertScheduledAlert(t *testing.T, got, want scheduledAlert) {
	t.Helper()
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

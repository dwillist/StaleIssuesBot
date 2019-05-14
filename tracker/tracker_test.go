package tracker_test

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/dwillist/stale_issues/tracker"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

//go:generate mockgen -destination=mocks_test.go -package=tracker_test github.com/dwillist/stale_issues/tracker Caller,Timer

func TestUnitTracker(t *testing.T) {
	spec.Run(t, "Tracker", testTracker, spec.Report(report.Terminal{}))
}

func testTracker(t *testing.T, when spec.G, it spec.S) {
	var (
		mockCtrl   *gomock.Controller
		mockCaller *MockCaller
		mockTimer  *MockTimer
		subject    tracker.Tracker
		response   string
		err        error
	)

	it.Before(func() {
		RegisterTestingT(t)
		mockCtrl = gomock.NewController(t)
		mockCaller = NewMockCaller(mockCtrl)
		mockTimer = NewMockTimer(mockCtrl)

		response, err = fileToString("tracker.json")
		Expect(err).ToNot(HaveOccurred())

		mockCaller.EXPECT().Call(tracker.Endpoint).Return(response, nil)

		now := time.Date(
			2019, 05, 14, 20, 34, 58, 651387237, time.UTC)
		mockTimer.EXPECT().Time().Return(now).AnyTimes()

		subject = tracker.NewTracker(mockCaller, mockTimer)
	})

	it.After(func() {
		mockCtrl.Finish()
	})

	when("Search()", func() {
		it("returns all open unstarted github issues", func() {
			result, err := subject.Search()
			Expect(err).ToNot(HaveOccurred())

			Expect(result[0].ID).To(Equal(165068541))
			Expect(len(result)).To(Equal(45))
		})
	})

	when("FilterIsues()", func() {
		it("returns all stale Github issues", func() {
			result, err := subject.Filter()
			Expect(err).ToNot(HaveOccurred())

			Expect(result).NotTo(BeEmpty())
			Expect(result[0].ID).NotTo(Equal(165068541))
			Expect(len(result)).To(Equal(16))
			Expect(result[0].ID).To(Equal(165092470))
		})
	})

}

func fileToString(fileName string) (string, error) {
	path, err := filepath.Abs(path.Join("..", "resources", "testdata", fileName))
	if err != nil {
		return "", err
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

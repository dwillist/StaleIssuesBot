package tracker_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"

	"github.com/dwillist/stale_issues/resources"

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
		response   []byte
		err        error
	)

	it.Before(func() {
		RegisterTestingT(t)
		mockCtrl = gomock.NewController(t)
		mockCaller = NewMockCaller(mockCtrl)
		mockTimer = NewMockTimer(mockCtrl)

		response, err = FileBytes("tracker.json")
		Expect(err).ToNot(HaveOccurred())

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
			mockCaller.EXPECT().Get(tracker.SearchEndpoint).Return(response, nil)
			result, err := subject.Search()
			Expect(err).ToNot(HaveOccurred())

			Expect(result[0].ID).To(Equal(165068541))
			Expect(len(result)).To(Equal(45))
		})
	})

	when("FilterIsues()", func() {
		it("returns all stale Github issues", func() {
			mockCaller.EXPECT().Get(tracker.SearchEndpoint).Return(response, nil)
			result, err := subject.Filter()
			Expect(err).ToNot(HaveOccurred())

			Expect(result).NotTo(BeEmpty())
			Expect(result[0].ID).NotTo(Equal(165068541))
			Expect(len(result)).To(Equal(16))
			Expect(result[0].ID).To(Equal(165092470))
		})
	})

	when("Initializing 'Stale' label", func() {
		it("doesn't create a label if one exists", func() {
			errorResponse, err := FileBytes("trackererror.json")
			Expect(err).NotTo(HaveOccurred())

			labelResponse, err := FileBytes("labels.json")

			labelStruct := resources.Label{Name: tracker.StaleLabel}
			labelJSON, err := json.Marshal(labelStruct)
			Expect(err).NotTo(HaveOccurred())

			mockCaller.EXPECT().Post(tracker.LabelsEndpoint, labelJSON).Return(errorResponse, nil)
			mockCaller.EXPECT().Get(tracker.LabelsEndpoint).Return(labelResponse, nil)
			label, success, err := subject.PostLabel()
			Expect(err).NotTo(HaveOccurred())
			Expect(success).To(BeFalse())
			Expect(label.Name).To(Equal("stale-issue"))
		})
		it("Creates the label if it does not exist", func() {
			labelResponse, err := FileBytes("new_label.json")
			Expect(err).NotTo(HaveOccurred())

			labelStruct := resources.Label{Name: tracker.StaleLabel}
			labelJSON, err := json.Marshal(labelStruct)
			Expect(err).NotTo(HaveOccurred())

			mockCaller.EXPECT().Post(tracker.LabelsEndpoint, labelJSON).Return(labelResponse, nil)
			label, success, err := subject.PostLabel()
			Expect(err).NotTo(HaveOccurred())
			Expect(success).To(BeTrue())
			Expect(label.Name).To(Equal("stale-issue"))
		})

		it("applies stale-issue tag to old issues", func() {
			postUpdateResponse, err := FileBytes("post_update_story.json")
			Expect(err).NotTo(HaveOccurred())

			preUpdateData, err := FileBytes("pre_update_story.json")
			Expect(err).NotTo(HaveOccurred())

			var preUpdateStory resources.Story

			Expect(json.Unmarshal(preUpdateData, &preUpdateStory)).To(Succeed())

			staleLabel := resources.Label{
				Name: tracker.StaleLabel,
			}

			newLabelData, err := json.Marshal(staleLabel)
			Expect(err).ToNot(HaveOccurred())

			// what is the data we are posting
			mockCaller.EXPECT().Post(tracker.StoriesEndpoint, newLabelData).Return(postUpdateResponse, nil)
			updatedStory, success, err := subject.UpdateStory(preUpdateStory)

			Expect(err).NotTo(HaveOccurred())
			Expect(success).To(BeTrue())
			Expect(updatedStory.Labels[1].Name).To(Equal(tracker.StaleLabel))
		})

		it("fails to re-apply the stale issues tag", func() {
			startData, err := FileBytes("post_update_story.json")
			Expect(err).NotTo(HaveOccurred())

			var startStory resources.Story
			Expect(json.Unmarshal(startData, &startStory)).To(Succeed())

			// what is the data we are posting
			updatedStory, success, err := subject.UpdateStory(startStory)

			Expect(err).NotTo(HaveOccurred())
			Expect(success).To(BeFalse())
			Expect(updatedStory).To(Equal(startStory))

		})

	})

}

// TODO
// Apply the tag to an issue

func FileBytes(fileName string) ([]byte, error) {
	path, err := filepath.Abs(filepath.Join("..", "resources", "testdata", fileName))
	if err != nil {
		return []byte{}, err
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return buf, nil
}

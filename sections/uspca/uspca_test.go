package uspca

import (
	"testing"

	"github.com/prebid/go-gpp/constants"
	"github.com/prebid/go-gpp/sections"
	"github.com/stretchr/testify/assert"
)

type uspcaTestData struct {
	description string
	gppString   string
	expected    USPCA
}

func TestUSPCA(t *testing.T) {
	testData := []uspcaTestData{
		{
			description: "should populate USPCA segments correctly",
			gppString:   "xlgWEYCY",
			/*
				110001 10 01 01 10 00 000101100001000110 0000 00 10 01 10
			*/
			expected: USPCA{
				CoreSegment: sections.USPCACoreSegment{
					Version:                     49,
					SalesOptOutNotice:           2,
					SharingOptOutNotice:         1,
					SensitiveDataLimitUseNotice: 1,
					SalesOptOut:                 2,
					SharingOptOut:               0,
					SensitiveDataProcessing: []byte{
						0, 1, 1, 2, 0, 1, 0, 1, 2,
					},
					KnownChildSensitiveDataConsents: []byte{
						0, 0,
					},
					PersonalDataConsents:    0,
					MspaCoveredTransaction:  2,
					MspaOptOutOptionMode:    1,
					MspaServiceProviderMode: 2,
				},
				GPCSegment: sections.CommonUSGPCSegment{
					Gpc: 0,
				},
				SectionID: constants.SectionUSPCA,
				Value:     "xlgWEYCY",
			},
		},
	}

	for _, test := range testData {
		result, err := NewUSPCA(test.gppString)

		assert.Nil(t, err)
		assert.Equal(t, test.expected, result)
		assert.Equal(t, constants.SectionUSPCA, result.GetID())
		assert.Equal(t, test.gppString, result.GetValue())
	}
}

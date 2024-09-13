package types_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cf "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cfTypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/michimani/cfkvs/types"
	"github.com/stretchr/testify/assert"
)

func Test_KVS_Parse(t *testing.T) {
	cases := []struct {
		name    string
		k       *types.KVS
		o       any
		expect  types.KVS
		wantErr bool
	}{
		{
			name: "CloudFront CreateKeyValueStoreOutput",
			k:    &types.KVS{},
			o: &cf.CreateKeyValueStoreOutput{
				KeyValueStore: &cfTypes.KeyValueStore{
					Id:      aws.String("id"),
					Name:    aws.String("name"),
					Comment: aws.String("comment"),
					Status:  aws.String("status"),
					ARN:     aws.String("arn"),
				},
			},
			expect: types.KVS{
				Id:      "id",
				Name:    "name",
				Comment: "comment",
				Status:  "status",
				ARN:     "arn",
			},
		},
		{
			name: "CloudFront DescribeKeyValueStoreOutput",
			k:    &types.KVS{},
			o: &cf.DescribeKeyValueStoreOutput{
				KeyValueStore: &cfTypes.KeyValueStore{
					Id:      aws.String("id"),
					Name:    aws.String("name"),
					Comment: aws.String("comment"),
					Status:  aws.String("status"),
					ARN:     aws.String("arn"),
				},
			},
			expect: types.KVS{
				Id:      "id",
				Name:    "name",
				Comment: "comment",
				Status:  "status",
				ARN:     "arn",
			},
		},
		{
			name: "CloudFront KeyValueStore",
			k:    &types.KVS{},
			o: &cfTypes.KeyValueStore{
				Id:      aws.String("id"),
				Name:    aws.String("name"),
				Comment: aws.String("comment"),
				Status:  aws.String("status"),
				ARN:     aws.String("arn"),
			},
			expect: types.KVS{
				Id:      "id",
				Name:    "name",
				Comment: "comment",
				Status:  "status",
				ARN:     "arn",
			},
		},
		{
			name: "nil KVS",
			k:    nil,
			o: &cf.CreateKeyValueStoreOutput{
				KeyValueStore: &cfTypes.KeyValueStore{
					Id:      aws.String("id"),
					Name:    aws.String("name"),
					Comment: aws.String("comment"),
					Status:  aws.String("status"),
					ARN:     aws.String("arn"),
				},
			},
			wantErr: true,
		},
		{
			name:    "invalid type",
			k:       &types.KVS{},
			o:       "invalid",
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.k.Parse(c.o)
			if c.wantErr {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, *c.k)
		})
	}
}

func Test_KVSList_Parse(t *testing.T) {
	cases := []struct {
		name    string
		kl      *types.KVSList
		o       *cf.ListKeyValueStoresOutput
		expect  types.KVSList
		wantErr bool
	}{
		{
			name: "CloudFront ListKeyValueStoresOutput",
			kl:   &types.KVSList{},
			o: &cf.ListKeyValueStoresOutput{
				KeyValueStoreList: &cfTypes.KeyValueStoreList{
					Items: []cfTypes.KeyValueStore{
						{
							Id:      aws.String("id1"),
							Name:    aws.String("name1"),
							Comment: aws.String("comment1"),
							Status:  aws.String("status1"),
							ARN:     aws.String("arn1"),
						},
						{
							Id:      aws.String("id2"),
							Name:    aws.String("name2"),
							Comment: aws.String("comment2"),
							Status:  aws.String("status2"),
							ARN:     aws.String("arn2"),
						},
					},
				},
			},
			expect: types.KVSList{
				{
					Id:      "id1",
					Name:    "name1",
					Comment: "comment1",
					Status:  "status1",
					ARN:     "arn1",
				},
				{
					Id:      "id2",
					Name:    "name2",
					Comment: "comment2",
					Status:  "status2",
					ARN:     "arn2",
				},
			},
		},
		{
			name: "nil KVSList",
			kl:   nil,
			o: &cf.ListKeyValueStoresOutput{
				KeyValueStoreList: &cfTypes.KeyValueStoreList{
					Items: []cfTypes.KeyValueStore{
						{
							Id:      aws.String("id1"),
							Name:    aws.String("name1"),
							Comment: aws.String("comment1"),
							Status:  aws.String("status1"),
							ARN:     aws.String("arn1"),
						},
						{
							Id:      aws.String("id2"),
							Name:    aws.String("name2"),
							Comment: aws.String("comment2"),
							Status:  aws.String("status2"),
							ARN:     aws.String("arn2"),
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(tt *testing.T) {
			asst := assert.New(tt)

			err := c.kl.Parse(c.o)
			if c.wantErr {
				asst.Error(err)
				return
			}

			asst.NoError(err)
			asst.Equal(c.expect, *c.kl)
		})
	}
}

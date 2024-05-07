package service_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github/Igo87/crypt/models"
	"github/Igo87/crypt/service/mocks"
)

func Test_AddCurrencies(t *testing.T) {
	t.Parallel()
	mockRepo := mocks.NewMockrepositoryManager(gomock.NewController(t))

	defer mockRepo.EXPECT().AddCurrencies(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	err := mockRepo.AddCurrencies(context.Background(), models.Crypto{})
	assert.NoError(t, err)
}

func Test_SelectCurrenciesByToday(t *testing.T) {
	t.Parallel()
	mockRepo := mocks.NewMockrepositoryManager(gomock.NewController(t))
	mockRepo.EXPECT().SelectCurrenciesByToday(gomock.Any()).Return(models.Currencies{}, nil).Times(1)

	ctx := context.TODO()

	result, err := mockRepo.SelectCurrenciesByToday(ctx)
	require.NoError(t, err)
	assert.Equal(t, models.Currencies{}, result)
}

func Test_SelectCurrenciesByLastDay(t *testing.T) {
	t.Parallel()
	mockRepo := mocks.NewMockrepositoryManager(gomock.NewController(t))
	mockRepo.EXPECT().SelectCurrenciesByLastDay(gomock.Any()).Return(models.Currencies{}, nil).Times(1)

	ctx := context.TODO()

	result, err := mockRepo.SelectCurrenciesByLastDay(ctx)
	require.NoError(t, err)
	assert.Equal(t, models.Currencies{}, result)
}

func Test_SelectCurrenciesByName(t *testing.T) {
	mockRepo := mocks.NewMockrepositoryManager(gomock.NewController(t))

	ctx := context.TODO()
	testCases := []struct {
		name string
	}{
		{"Bitcoin"},
		{"Ethereum"},
	}
	expected := models.Currencies{
		{
			Name:    "Bitcoin",
			Percent: 1.0,
			Min:     99.0,
			Max:     100.0,
		},
		{
			Name:    "Ethereum",
			Percent: 10.0,
			Min:     90.0,
			Max:     100.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.EXPECT().SelectCurrenciesByName(gomock.Any(), tc.name).Return(expected, nil).Times(1)
			res, err := mockRepo.SelectCurrenciesByName(ctx, tc.name)
			assert.Equal(t, expected, res)
			assert.NoError(t, err)
		})
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockrepositoryManager(ctrl)
	mockRepo.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(models.Crypto{}, nil).Times(1)
	res, err := mockRepo.Fetch(context.Background(), slog.Logger{})
	require.NoError(t, err)
	assert.Equal(t, models.Crypto{}, res)

}
func TestRun(t *testing.T) {
	t.Parallel()
	mockRepo := mocks.NewMockrepositoryManager(gomock.NewController(t))
	mockRepo.EXPECT().Run(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	err := mockRepo.Run(context.Background(), slog.Logger{})
	assert.NoError(t, err)
}

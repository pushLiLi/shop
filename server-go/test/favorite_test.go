package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type FavoriteTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (s *FavoriteTestSuite) SetupSuite() {
	SetupTestConfig()
	SetupTestDB()
	s.router = SetupRouter()
}

func TestFavoriteSuite(t *testing.T) {
	suite.Run(t, new(FavoriteTestSuite))
}

func (s *FavoriteTestSuite) TestGetFavorites() {
	w := MakeRequest(s.router, "GET", "/api/favorites", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	s.Equal(3, len(items), "seeded user should have 3 favorites")
}

func (s *FavoriteTestSuite) TestGetFavoritesNotLoggedIn() {
	w := MakeRequest(s.router, "GET", "/api/favorites", nil, nil)
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	s.Equal(0, len(items))
}

func (s *FavoriteTestSuite) TestAddFavorite() {
	productID := Data.Products[3].ID
	body := map[string]interface{}{
		"productId": productID,
	}
	w := MakeRequest(s.router, "POST", "/api/favorites", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
}

func (s *FavoriteTestSuite) TestAddFavoriteDuplicate() {
	productID := Data.Products[0].ID
	body := map[string]interface{}{
		"productId": productID,
	}
	w := MakeRequest(s.router, "POST", "/api/favorites", body, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	s.Equal(true, resp["exists"])
}

func (s *FavoriteTestSuite) TestAddFavoriteUnauthorized() {
	body := map[string]interface{}{
		"productId": 1,
	}
	w := MakeRequest(s.router, "POST", "/api/favorites", body, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *FavoriteTestSuite) TestDeleteFavorite() {
	productID := Data.Products[0].ID
	w := MakeRequest(s.router, "DELETE", fmt.Sprintf("/api/favorites/%d", productID), nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)

	w = MakeRequest(s.router, "GET", "/api/favorites", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["items"].([]interface{})
	s.Equal(3, len(items), "3 original favorites + 1 added by TestAddFavorite - 1 deleted = 3")
}

func (s *FavoriteTestSuite) TestDeleteFavoriteUnauthorized() {
	w := MakeRequest(s.router, "DELETE", "/api/favorites/1", nil, nil)
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *FavoriteTestSuite) TestDeleteFavoriteInvalidProductID() {
	w := MakeRequest(s.router, "DELETE", "/api/favorites/abc", nil, GetCustomerAuthHeader())
	s.Equal(http.StatusBadRequest, w.Code)
}

package controller

import (
	"app/internal/usecase/repo"
	"app/pkg/httpserver"
	"app/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct {
	routs []httpserver.Route
}

func NewController() *Controller {
	logger.Debug("Controller initialized")
	return &Controller{}
}

// GetRoutes returns the list of all the routes that this handler exposes
func (lc *Controller) GetRoutes() []httpserver.Route {
	routes := make([]httpserver.Route, 0)

	routes = append(routes, httpserver.Route{Name: "Get users balance by id", Method: http.MethodGet, Pattern: "/{userId}",
		HandlerFunc: lc.getBalanceByUserId})

	routes = append(routes, httpserver.Route{Name: "Get history of transanctions of user by id", Method: http.MethodGet, Pattern: "/history/{userId}",
		HandlerFunc: lc.getInfoTransanctions})

	routes = append(routes, httpserver.Route{Name: "Change balance", Method: http.MethodPost, Pattern: "/{userId}",
		HandlerFunc: lc.increaseOrDecreaseBalance})

	routes = append(routes, httpserver.Route{Name: "Transanction between users", Method: http.MethodPost, Pattern: "/{userId1}/{userId2}",
		HandlerFunc: lc.transanct})
	return routes
}

func (lc *Controller) getBalanceByUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["userId"], 10, 0)
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid user_id"))
	}

	user, err := repo.GetUserBalance(int(userId))
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "User with this Id does not exist"))
	}

	httpserver.WriteResponse(w, http.StatusOK, user)
}

func (lc *Controller) getInfoTransanctions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["userId"], 10, 0)
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid user_id"))
	}
	history, err := repo.GetUserTransanctions(int(userId))
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "User with this Id does not exist"))
	}
	httpserver.WriteResponse(w, http.StatusOK, history)
}

func (lc *Controller) increaseOrDecreaseBalance(w http.ResponseWriter, r *http.Request) {
	var am float64
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid user_id"))
	}
	in := r.URL.Query().Get("incr")
	de := r.URL.Query().Get("decr")
	if len(in) == 0 && len(de) == 0 {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Missing action"))
	}

	if len(in) > 0 {
		am, err = strconv.ParseFloat(in, 64)
	} else {
		am, err = strconv.ParseFloat(de, 64)
		am = -am
	}

	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid value for operation"))
	}

	err = repo.ChangeBalance(userId, am)

	if err != nil {
		logger.Err("Err of method ChangeBalance", err)
	} else {
		httpserver.WriteResponse(w, http.StatusOK, "Balance operation was successful")
	}
}

func (lc *Controller) transanct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId1, err := strconv.ParseInt(vars["userId1"], 10, 0)
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid user_id sender"))
	}
	userId2, err := strconv.ParseInt(vars["userId2"], 10, 0)
	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid user_id recipient"))
	}

	sum := r.URL.Query().Get("sum")

	am, err := strconv.ParseFloat(sum, 64)

	if err != nil {
		httpserver.WriteResponse(w, http.StatusBadRequest, httpserver.NewError(httpserver.BadRequestError, "Invalid value for operation"))
	}

	err = repo.TransanctionBetweenUsers(int(userId1), int(userId2), am)

	if err != nil {
		logger.Err("Err of method ChangeBalance", err)
		httpserver.WriteResponse(w, http.StatusOK, "Transanction between users was denied")
	} else {
		httpserver.WriteResponse(w, http.StatusOK, "Transanction between users was successful")
	}
}

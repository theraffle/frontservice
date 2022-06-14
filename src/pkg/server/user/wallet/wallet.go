/*
 Copyright 2022 The Raffle Authors

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package wallet

import (
	"github.com/go-logr/logr"
	"github.com/theraffle/frontservice/src/internal/apiserver"
	"github.com/theraffle/frontservice/src/internal/wrapper"
	"net/http"
)

type handler struct {
	_ logr.Logger
}

// NewHandler instantiates a new apis handler
func NewHandler(parent wrapper.RouterWrapper, _ logr.Logger) (apiserver.APIHandler, error) {
	handler := &handler{}

	// Create User Project
	createUserWallet := wrapper.New("/wallet", []string{http.MethodPost}, handler.createUserWalletHandler)
	if err := parent.Add(createUserWallet); err != nil {
		return nil, err
	}

	// Get User Projects
	getUserWallet := wrapper.New("/wallet", []string{http.MethodGet}, handler.getUserWalletHandler)
	if err := parent.Add(getUserWallet); err != nil {
		return nil, err
	}

	return handler, nil
}

func (h handler) getUserWalletHandler(writer http.ResponseWriter, request *http.Request) {

}

func (h handler) createUserWalletHandler(writer http.ResponseWriter, request *http.Request) {

}

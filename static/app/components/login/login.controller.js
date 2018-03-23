'use strict';

var app = angular.module('ipl');

/**
 * Login Controller
 * 
 * Controller for the login page.
 */
app.controller('loginController', function ($http, $window, $state, INumberPattern, token, urlService) {
    var vm = this;

    vm.iNumberPattern = INumberPattern;

    vm.signIn = signIn;

    // Function when sign in occurs
    function signIn() {
        vm.error = false;
        var data = {
            inumber: vm.iNumber,
            password: vm.password
        };
        $http.post(urlService.loginUser, data)
            .then(function (res) {
                console.log('login success',res.data,res.body);
                $window.localStorage.setItem(token, res.data.token);
                // Add JWT Token as the default token for all back-end requests
                $http.defaults.headers.common.Authorization = res.data.token;
                $state.go('main.home');

            }, function (err) {
                console.log('error', err);
                vm.error = true;
            });
    }
});

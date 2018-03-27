'use strict';

var app = angular.module('ipl');

/**
 * Login Controller
 * 
 * Controller for the login page.
 */
app.controller('loginController', function ($http, $window, $state, INumberPattern, urlService, utilsService) {
    var vm = this;

    vm.iNumberPattern = INumberPattern;

    vm.signIn = signIn;

    // Function when sign in occurs
    function signIn(isFormValid) {
        if (isFormValid === false) {
            utilsService.showToast({
                text: 'Please enter valid credentials.',
                hideDelay: 0,
                isError: true
            });
            return;
        }
        var data = {
            inumber: vm.iNumber,
            password: vm.password
        };
        var params = {
            url: urlService.loginUser,
            data: data,
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        };
        $http(params)
            .then(function (res) {
                utilsService.showToast({
                    text: 'Login Successful.',
                    hideDelay: 3000,
                    isError: false
                });
                console.log('login success');
                $window.localStorage.setItem('token', res.data.token);
                $window.localStorage.setItem('iNumber', vm.iNumber);
                // Add JWT Token as the default token for all back-end requests
                $http.defaults.headers.common.Authorization = res.data.token;
                $state.go('main.profile');
            }, function (err) {
                console.log('error', err);
                utilsService.showToast({
                    text: 'Please check your credentials.',
                    hideDelay: 0,
                    isError: true
                });
            });
    }
});

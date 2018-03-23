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
<<<<<<< 84e806b8c71d995b0992b8de37324f2ff2950330
                console.log('login success',res.data,res.body);
                $window.localStorage.setItem(token, res.data.token);
=======
                utilsService.showToast({
                    text: 'Login Successful.',
                    hideDelay: 3000,
                    isError: false
                });
                console.log('login success');
                $window.localStorage.setItem('token', res.data.token);
                $window.localStorage.setItem('iNumber', vm.iNumber);
>>>>>>> Profile and edit profile changes
                // Add JWT Token as the default token for all back-end requests
                $http.defaults.headers.common.Authorization = res.data.token;
                $state.go('main.home');
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

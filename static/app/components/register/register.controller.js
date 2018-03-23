'use strict';

var app = angular.module('ipl');

/**
 * Register Controller
 * 
 * Controller for the register page.
 */
app.controller('registerController', function ($http, INumberPattern, urlService, aliasPattern, utilsService) {
    var vm = this;

    vm.iNumberPattern = INumberPattern;
    vm.aliasPattern = aliasPattern;

    vm.signUp = signUp;

    // Function to sign up new user
    function signUp(isFormValid) {
        if (isFormValid === false) {
            utilsService.showToast({
                text: 'Please enter valid credentials.',
                hideDelay: 3000,
                isError: true
            });
            return;
        }
        if (vm.password !== vm.confirmPassword) {
            utilsService.showToast({
                text: 'Password and Confirm Password do not match',
                hideDelay: 3000,
                isError: true
            });
            return;
        }
        var data = {
            iNumber: vm.iNumber,
            firstName: vm.firstName,
            lastName: vm.lastName,
            password: vm.password
        };
        if (vm.alias !== '' && vm.alias !== undefined && vm.alias !== null) {
            data.alias = vm.alias;
        }
        var params = {
            url: urlService.registerUser,
            data: data,
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        };
        $http(params)
            .then(function () {
                console.log('signup success');
                utilsService.showToast({
                    text: 'User Registration Successful.',
                    hideDelay: 0,
                    isError: true
                });
            }, function (err) {
                console.log('signup error', err);
                utilsService.showToast({
                    text: `${err.message} .`,
                    hideDelay: 0,
                    isError: true
                });
            });
    }

});
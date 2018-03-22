'use strict';

var app = angular.module('ipl');

/**
 * Register Controller
 * 
 * Controller for the register page.
 */
app.controller('registerController', function ($http, INumberPattern, urlService) {
    var vm = this;

    vm.iNumberPattern = INumberPattern;
    vm.successMessage = 'User registration successful';

    vm.signUp = signUp;

    // Function to sign up new user
    function signUp() {
        vm.error = false;
        vm.errorMessage = undefined;
        vm.success = false;
        if (vm.password !== vm.confirmPassword) {
            vm.success = false;
            vm.error = true;
            vm.errorMessage = 'Error! Password and Confirm Password do not match.';
            return;
        }
        var data = {
            iNumber: vm.iNumber,
            firstName: vm.firstName,
            lastName: vm.lastName,
            password: vm.password
        };
        vm.error = false;
        if (vm.alias !== '' && vm.alias !== undefined && vm.alias !== null) {
            data.alias = vm.alias;
        }
        $http.post(urlService.registerUser, data)
            .then(function () {
                console.log('signup success');
                vm.success = true;
                vm.error = false;
            }, function (err) {
                console.log('signup error', err);
                vm.error = true;
                vm.success = false;
                vm.errorMessage = err;
            });
    }

});
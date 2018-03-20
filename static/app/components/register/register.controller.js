'use strict';

var app = angular.module('ipl');

/**
 * Register Controller
 * 
 * Controller for the register page.
 */
app.controller('registerController', function (requestsService) {
    var vm = this;

    vm.abcd = 'register';

    vm.signUp = signUp;

    function signUp() {
        var data = {
            iNumber: vm.iNumber,
            firstName: vm.firstName,
            lastName: vm.lastName,
            password: vm.password
        };
        if (vm.alias !== '' && vm.alias !== undefined) {
            data.alias = vm.alias;
        }
        requestsService.registerUser(data)
            .then(function () {
                console.log('signup success');
            })
            .catch(function (err) {
                console.log('signup error', err);
            });
    }
});
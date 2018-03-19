'use strict';

var app = angular.module('ipl');

/**
 * Login Controller
 * 
 * Controller for the login page.
 */
app.controller('loginController', function () {
    var vm = this;

    vm.abcd = 'login';
    vm.iNumberPattern = /^[i|I][0-9]{6}$/;
});
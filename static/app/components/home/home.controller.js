'use strict';

var app = angular.module('ipl');

app.controller('homeController', function ($window) {
    var vm = this;

    vm.abcd = 'akshil';
    // $window.localStorage.clear();
});
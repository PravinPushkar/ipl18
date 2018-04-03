'use strict';

var app = angular.module('ipl');

/**
 * Rules Controller
 * 
 * Controller for rules page
 */
app.controller('rulesController', ['rulesService', function (rulesService) {
    var vm = this;

    vm.rulesList = rulesService;
}]);
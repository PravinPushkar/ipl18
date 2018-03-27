'use strict';

var app = angular.module('ipl');

/**
 * Bottom Sheet Controller
 * 
 * Controller for the bottom sheet displayed
 * instead of the user menu in small screens.
 */
app.controller('bottomSheetGridController', function(toolbarService) {
    var vm = this;

    vm.menuItems = toolbarService.userMenuItems;
});
'use strict';

var app = angular.module('ipl');

/**
 * Bottom Sheet Controller
 * 
 * Controller for the bottom sheet displayed
 * instead of the user menu in small screens.
 */
app.controller('bottomSheetGridController', ['$window', '$state', '$mdBottomSheet', 'utilsService', 'toolbarService', function ($window, $state, $mdBottomSheet, utilsService, toolbarService) {
    var vm = this;

    vm.menuItems = toolbarService.userMenuItems;
    vm.clickMenuItem = clickMenuItem;

    function clickMenuItem(id) {
        switch (id) {
        case 'profile':
            $state.go('main.profile');
            break;
        case 'logout':
            $window.localStorage.removeItem('token');
            $window.localStorage.removeItem('iNumber');
            utilsService.showToast({
                text: 'Logout Successful.',
                hideDelay: 1500,
                isError: false
            });
            $state.go('login');
            $mdBottomSheet.hide();
            break;
        default:
            utilsService.showToast({
                text: 'ID is not registered',
                hideDelay: 0,
                isError: true
            });
            break;
        }
    }
}]);
'use strict';

var app = angular.module('ipl');

/**
 * Toolbar Controller
 * 
 * Controller for the toolbar and the sidebar.
 */
app.controller('toolbarController', ['$mdSidenav', '$mdComponentRegistry', '$rootScope', '$http', '$window', '$state', '$mdBottomSheet', 'toolbarService', 'utilsService', function ($mdSidenav, $mdComponentRegistry, $rootScope, $http, $window, $state, $mdBottomSheet, toolbarService, utilsService) {
    var vm = this;

    vm.toggleSidenav = toggleSidenav;
    vm.clickUserMenu = clickUserMenu;

    vm.sidenavId = 'left';
    vm.sidebarItems = toolbarService.sidebarItems;
    vm.userMenuItems = toolbarService.userMenuItems;

    // Change the profile pic in the toolbar by watching the localstorage
    function profilePicSet() {
        return $window.localStorage.getItem('profilePic');
    }

    $rootScope.$watch(profilePicSet, function (picLocation) {
        vm.imageStyle = {
            'background-image': `url('${picLocation}')`,
            'background-size': 'cover',
            'background-position': 'center center'
        };
    });

    $rootScope.$on('$locationChangeStart', function () {
        $mdSidenav(vm.sidenavId).close();
    });

    // Display bottom sheet in moble view
    vm.showGridBottomSheet = function () {
        $mdBottomSheet.show({
            templateUrl: '/static/app/shared/toolbar/bottomSheetGrid.html',
            controller: 'bottomSheetGridController',
            controllerAs: 'bottomSheet',
            clickOutsideToClose: true
        });
    };

    // Function to toggle the sidebar visibility
    function toggleSidenav() {
        $mdComponentRegistry
            .when(vm.sidenavId)
            .then(function () {
                $mdSidenav(vm.sidenavId, true).toggle();
            });
    }

    // Function for when user clicks on the user menu
    function clickUserMenu(id) {
        switch (id) {
        case 'profile':
            $state.go('main.profile');
            break;
        case 'logout':
            var params = {
                title: 'Confirm Logout',
                text: 'Are you sure you want to Logout?',
                aria: 'logout',
                ok: 'Yes',
                cancel: 'No'
            };
            utilsService.showConfirmDialog(params)
                .then(function () {
                    utilsService.logout('Logout Successful', false);
                }, function () {
                    console.log('Logout cancelled');
                });
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
'use strict';

var app = angular.module('ipl');

/**
 * Toolbar Controller
 * 
 * Controller for the toolbar and the sidebar.
 */
app.controller('toolbarController', function ($mdSidenav, $mdComponentRegistry, $rootScope, $http, $window, $state, displayName, token, toolbarService, utilsService, urlService) {
    var vm = this;

    vm.toggleSidenav = toggleSidenav;
    vm.clickUserMenu = clickUserMenu;

    vm.sidenavId = 'left';
    vm.sidebarItems = toolbarService.sidebarItems;
    vm.userMenuItems = toolbarService.userMenuItems;

    $rootScope.$on('$locationChangeStart', function () {
        $mdSidenav(vm.sidenavId).close();
    });

    // Function to toggle the sidebar visibility
    function toggleSidenav() {
        $mdComponentRegistry
            .when(vm.sidenavId)
            .then(function () {
                $mdSidenav(vm.sidenavId, true).toggle();
            });
    }

    // Function for when user clicks on the user menu
    function clickUserMenu(id, event) {
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
                cancel: 'No',
                event: event
            };
            utilsService.showConfirmDialog(params)
                .then(function () {
                    $window.localStorage.removeItem(displayName);
                    $window.localStorage.removeItem(token);
                    $http.delete(urlService.logoutUser)
                        .then(function () {
                            $http.defaults.headers.common.Authorization = '';
                            $state.go('login');
                        }, function (err) {
                            console.log('Error logging out', err.message);
                        });
                }, function () {
                    console.log('Logout cancelled');
                });
            break;
        default:
            console.log('Error, ID is not registered');
            break;
        }
    }
});
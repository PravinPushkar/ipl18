'use strict';

var app = angular.module('ipl');

/**
 * Toolbar Controller
 * 
 * Controller for the toolbar and the sidebar.
 */
app.controller('toolbarController', function ($mdSidenav, $mdComponentRegistry, $rootScope, $http, $window, $state, $mdBottomSheet, toolbarService, utilsService, urlService) {
    var vm = this;

    var token;
    
    vm.toggleSidenav = toggleSidenav;
    vm.clickUserMenu = clickUserMenu;

    vm.sidenavId = 'left';
    vm.sidebarItems = toolbarService.sidebarItems;
    vm.userMenuItems = toolbarService.userMenuItems;
    // akshil check this
    $rootScope.$on('$locationChangeSuccess', function (newState, oldState) {
        console.log('tt', newState, oldState);
        vm.profilePic = $window.localStorage.getItem('picLocation'); 
        vm.imageStyle = {
            'background-image': `url('${vm.profilePic}')`,
            'background-size': 'cover',
            'background-position': 'center center'
        };
        console.log('picloc', vm.profilePic);
    });

    $rootScope.$on('$locationChangeStart', function () {
        $mdSidenav(vm.sidenavId).close();
    });

    // Display bottom sheet in moble view
    vm.showGridBottomSheet = function () {
        vm.alert = '';
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
                    token = $window.localStorage.getItem('token');
                    var params = {
                        url: urlService.logoutUser,
                        method: 'DELETE',
                        headers: {
                            'Authorization': token
                        }
                    };
                    $http(params)
                        .then(function () {
                            $window.localStorage.removeItem('token');
                            $window.localStorage.removeItem('iNumber');
                            $window.localStorage.removeItem('picLocation');
                            $state.go('login');
                        }, function (err) {
                            $window.localStorage.removeItem('token');
                            $window.localStorage.removeItem('iNumber');
                            $window.localStorage.removeItem('picLocation');
                            $state.go('login');
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
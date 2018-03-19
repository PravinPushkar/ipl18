'use strict';

var app = angular.module('ipl');

/**
 * Toolbar Controller
 * 
 * Controller for the toolbar and the sidebar.
 */
app.controller('toolbarController', function ($mdSidenav, $mdComponentRegistry) {
    var vm = this;

    vm.toggleSidenav = toggleSidenav;

    // Array of items in the sidebar menu
    vm.sidebarItems = [{
        name: 'Home',
        icon: 'home'
    }, {
        name: 'Leaderboard',
        icon: 'assessment'
    }, {
        name: 'Rules',
        icon: 'assignment'
    }, {
        name: 'Teams',
        icon: 'people'
    }];

    // Array of items in the user menu
    vm.userMenuItems = [{
            name: 'Profile',
            icon: 'account_box'
        },
        {
            name: 'Logout',
            icon: 'exit_to_app'
        }
    ];

    // Function to toggle the sidebar visibility
    function toggleSidenav(sidenavId) {
        $mdComponentRegistry
            .when(sidenavId)
            .then(function () {
                $mdSidenav(sidenavId, true).toggle();
            });
    }
});
'use strict';

var app = angular.module('ipl');

/**
 * Toolbar service
 * 
 * This service contains menu items
 * for the toolbar and sidebar
 */
app.factory('toolbarService', function () {
    var service = {};

    service.sidebarItems = [{
            name: 'Home',
            icon: 'home',
            state: 'main.home'
        },
        {
            name: 'Leaderboard',
            icon: 'assessment',
            state: 'main.leaderboard'
        },
        {
            name: 'Rules',
            icon: 'assignment',
            state: 'main.rules'
        },
        {
            name: 'Teams',
            icon: 'people',
            state: 'main.teams'
        }
    ];
    service.userMenuItems = [{
            name: 'Profile',
            icon: 'account_box',
            id: 'profile'
        },
        {
            name: 'Logout',
            icon: 'exit_to_app',
            id: 'logout'
        }
    ];

    return service;
});
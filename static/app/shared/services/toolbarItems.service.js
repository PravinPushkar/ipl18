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
            name: 'Fixtures',
            icon: 'event',
            state: 'main.fixtures'
        },
        {
            name: 'Main Event Predictions',
            icon: 'star',
            state: 'main.mainEvent'
        },
        {
            name: 'Leaderboard',
            icon: 'assessment',
            state: 'main.leaderboard'
        },
        {
            name: 'Teams',
            icon: 'people',
            state: 'main.teams'
        },
        {
            name: 'Rules',
            icon: 'assignment',
            state: 'main.rules'
        },
        {
            name: 'Feeds',
            icon: 'people',
            state: 'main.feeds'
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
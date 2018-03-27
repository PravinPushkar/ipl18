'use strict';

var app = angular.module('ipl');

/**
 * Routing
 * 
 * The routing config for this application.
 * It uses states for routing purposes.
 */
app.config(function ($stateProvider, $urlRouterProvider, $locationProvider, $urlMatcherFactoryProvider) {

    // prefixing hash with '' to avoid hashbang
    $locationProvider.hashPrefix('');

    $urlMatcherFactoryProvider.caseInsensitive(true);

    // Array of state definitions, add additional states here
    var states = [{
        name: 'login',
        url: '/login',
        templateUrl: '/static/app/components/login/login.html',
        controller: 'loginController',
        controllerAs: 'login'
    }, {
        name: 'register',
        url: '/register',
        templateUrl: '/static/app/components/register/register.html',
        controller: 'registerController',
        controllerAs: 'register'
    }, {
        abstract: true,
        name: 'main',
        views: {
            '@': {
                templateUrl: '/static/app/components/main/main.html'
            },
            'top@main': {
                templateUrl: '/static/app/shared/toolbar/toolbar.html',
                controller: 'toolbarController',
                controllerAs: 'toolbar'
            }
        }
    }, {
        name: 'main.teams',
        url: '/teams',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/teams/teams.html',
                controller: 'teamsController',
                controllerAs: 'teams'
            }
        }
    }, {
        name: 'main.profile',
        url: '/profile',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/profile/profile.html',
                controller: 'profileController',
                controllerAs: 'profile'
            }
        }
    }, {
        name: 'main.editProfile',
        url: '/editprofile',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/editProfile/editProfile.html',
                controller: 'editProfileController',
                controllerAs: 'editProfile'
            }
        }
    }, {
        name: 'main.leaderboard',
        url: '/leaderboard',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/leaderboard/leaderboard.html',
                controller: 'leaderboardController',
                controllerAs: 'leaderboard'
            }
        }
    }];

    // Add every state into the $stateProvider
    states.forEach(function (state) {
        $stateProvider.state(state);
    });

    // Default page
    $urlRouterProvider.otherwise('/profile');
});
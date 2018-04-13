'use strict';

var app = angular.module('ipl');

/**
 * Routing
 * 
 * The routing config for this application.
 * It uses states for routing purposes.
 */
app.config(['$stateProvider', '$urlRouterProvider', '$locationProvider', '$urlMatcherFactoryProvider', function ($stateProvider, $urlRouterProvider, $locationProvider, $urlMatcherFactoryProvider) {

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
    }, {
        name: 'main.teams.players',
        url: '/teams/:teamId',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/players/players.html',
                controller: 'playersController',
                controllerAs: 'players'
            }
        }
    }, {
        name: 'main.fixtures',
        url: '/fixtures',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/fixtures/fixtures.html',
                controller: 'fixturesController',
                controllerAs: 'fixtures'
            }
        }
    }, {
        name: 'main.rules',
        url: '/rules',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/rules/rules.html',
                controller: 'rulesController',
                controllerAs: 'rules'
            }
        }
    }, {
        name: 'main.mainEvent',
        url: '/mainevent',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/mainEvent/mainEvent.html',
                controller: 'mainEventController',
                controllerAs: 'mainEvent'
            }
        }
    }, {
        name: 'main.feeds',
        url: '/feeds',
        views: {
            'body@main': {
                templateUrl: '/static/app/components/feeds/feeds.html',
                controller: 'feedsController',
                controllerAs: 'feeds'
            }
        }
    }];

    // Add every state into the $stateProvider
    states.forEach(function (state) {
        $stateProvider.state(state);
    });

    // Default page
    $urlRouterProvider.otherwise('/profile');
}]);
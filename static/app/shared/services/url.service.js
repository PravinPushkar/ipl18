'use strict';

var app = angular.module('ipl');

/**
 * URL service
 * 
 * This service contains the URL endpoints
 * of the back-end
 */
app.factory('urlService', function () {
    var service = {};

    service.registerUser = '/pub/register';
    service.loginUser = '/pub/login';
    service.logoutUser = '/api/logout';
    service.userProfile = '/api/users';
    service.leaderboard = '/api/leaders';
    service.teams = '/api/teams';
    service.bonus = '/api/bonus';
    service.players = '/api/players';
    service.fixtures = '/api/matches';
    service.predictions = '/api/predictions';
    service.feeds = '/feeds';

    return service;
});
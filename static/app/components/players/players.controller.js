'use strict';

var app = angular.module('ipl');

/**
 * Player controller
 * 
 * Controller for the player page.
 */
app.controller('playersController', ['$http', '$stateParams', '$window', 'urlService', 'utilsService', 'teamsDetailsService', function ($http, $stateParams, $window, urlService, utilsService, teamsDetailsService) {
    var vm = this;
    var token;

    vm.init = init;

    // Init function for players
    function init() {
        token = $window.localStorage.getItem('token');
        var playersParams = {
            url: `${urlService.teams}/${$stateParams.teamId}/players`,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        var teamParams = {
            url: `${urlService.teams}/${$stateParams.teamId}`,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        vm.playersList = [];
        var role;
        $http(teamParams)
            .then(function (res) {
                vm.teamAlias = res.data.shortName;
                vm.teamName = res.data.name;
                vm.teamDetails = teamsDetailsService[vm.teamAlias];
            }, function (err) {
                console.log('err', err);
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in fetching team details',
                    hideDelay: 0,
                    isError: true
                });
            });
        $http(playersParams)
            .then(function (res) {
                console.log('res', res);
                res.data.players.forEach(function (player) {
                    if (player.role === 'allrounder') {
                        role = 'All-Rounder';
                    } else {
                        role = utilsService.capitalizeFirstLetter(player.role);
                    }
                    vm.playersList.push({
                        playerId: player.id,
                        name: player.name,
                        role: role,
                        teamId: player.teamId
                    });
                });
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in fetching player details',
                    hideDelay: 0,
                    isError: true
                });
            });
    }
}]);
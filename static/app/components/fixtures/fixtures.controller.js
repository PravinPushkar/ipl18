'use strict';

var app = angular.module('ipl');

/**
 * Fixtures Controller
 * 
 * Controller for fixtures page.
 */
app.controller('fixturesController', ['$http', '$window', '$q', '$mdDialog', 'utilsService', 'urlService', function ($http, $window, $q, $mdDialog, utilsService, urlService) {
    var vm = this;
    var token;
    var iNumber;

    vm.init = init;
    vm.flag = [];
    vm.coinIcon = [];
    vm.coinFlag = [];
    vm.selectedTeam = [];
    vm.searchItem = '';
    vm.toggleFlag = toggleFlag;
    vm.useCoin = useCoin;
    vm.makePreditction = makePreditction;
    vm.clearSearchItem = clearSearchItem;
    vm.playerInTeam = playerInTeam;
    vm.openPredictionDialog = openPredictionDialog;

    // Finds if player is part of the playing teams
    function playerInTeam(playerTeamId, teamId1, teamId2) {
        // Akshil -> change this to === later
        if (playerTeamId == teamId1 || playerTeamId == teamId2) {
            return true;
        } else {
            return false;
        }
    }

    // Clears the search bar for select
    function clearSearchItem() {
        vm.searchItem = '';
    }

    // Toggle a flag and select the team chosen to predict
    function toggleFlag(id, teamId1, teamId2) {
        vm.flag[id][teamId1] = true;
        vm.flag[id][teamId2] = false;
        vm.selectedTeam[id] = teamId1;
    }

    // Check if user uses coin
    function useCoin(id) {
        vm.coinFlag[id] = !vm.coinFlag[id];
        if (vm.coinFlag[id] === false) {
            vm.coinIcon[id] = 'attach_money';
        } else {
            vm.coinIcon[id] = 'money_off';
        }
    }

    // Function to get team object from team id
    function getTeamFromId(teamId) {
        return vm.teamsList.find(function (team) {
            return team.id === teamId;
        });
    }

    // Init function for main fixtures view
    function init() {
        token = $window.localStorage.getItem('token');
        iNumber = $window.localStorage.getItem('iNumber');
        var fixturesParams = {
            url: urlService.fixtures,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        var teamParams = {
            url: urlService.teams,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        var playersParams = {
            url: urlService.players,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        vm.fixturesList = [];
        vm.teamsList = [];
        vm.playersList = [];
        var teamPromise = $http(teamParams);
        var playerPromise = $http(playersParams);
        // Resolve both promises
        $q.all([teamPromise, playerPromise])
            .then(function (data) {
                data[0].data.teams.forEach(function (team) {
                    vm.teamsList.push({
                        id: team.id,
                        name: team.name,
                        alias: team.shortName,
                        teamPic: team.picLocation
                    });
                });
                var role;
                data[1].data.players.forEach(function (player) {
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

                $http(fixturesParams)
                    .then(function (res) {
                        res.data.matches.forEach(function (fixture) {
                            vm.fixturesList.push({
                                teamId1: fixture.teamId1,
                                teamId2: fixture.teamId2,
                                venue: fixture.venue,
                                timestamp: moment(fixture.date).format('MMMM Do YYYY, h:mm:ss a'),
                                status: fixture.status,
                                matchId: fixture.id,
                                result: fixture.winner,
                                manOfMatch: fixture.mom,
                                star: fixture.star,
                                lockPred: fixture.lock,
                                team1: getTeamFromId(fixture.teamId1),
                                team2: getTeamFromId(fixture.teamId2),
                            });
                            vm.flag[fixture.mid] = {};
                            vm.coinIcon[fixture.id] = 'money_off';
                            vm.coinFlag[fixture.id] = false;
                        });
                        vm.fixturesList.sort(function(a, b){
                            return a.matchId-b.matchId
                        })

                    }, function (err) {
                        if (err.data.code === 403 && err.data.message === 'token not valid') {
                            utilsService.logout('Session expired, please re-login', true);
                            return;
                        }
                        utilsService.showToast({
                            text: 'Error in fetching fixtures',
                            hideDelay: 0,
                            isError: true
                        });
                    });

            })
            .catch(function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in fetching fixtures',
                    hideDelay: 0,
                    isError: true
                });
            });
    }

    // Function to send prediction data to the backend
    function makePreditction(id) {
        var data = {
            mid: id,
            inumber: iNumber,
            vote_team: vm.selectedTeam[id],
            vote_mom: vm.mom[id],
            coinused: vm.coinFlag[id]
        };
        var params = {
            url: urlService.predictions,
            method: 'POST',
            data: data,
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        $http(params)
            .then(function () {
                utilsService.showToast({
                    text: 'Successfully submitted prediction',
                    hideDelay: 1500,
                    isError: false
                });
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in submitting prediction, try again later',
                    hideDelay: 0,
                    isError: true
                });
            });
    }

    function openPredictionDialog(event, id) {
        $mdDialog.show({
            templateUrl: '/static/app/components/fixtures/fixturesDialog.html',
            controller: fixturesDialogController,
            targetEvent: event,
            locals: {
                matchId: id
            },
            clickOutesideToClose: true
        });
    }

    function fixturesDialogController($scope, matchId) {
        $scope.playersList = vm.playersList;
    }

}]);
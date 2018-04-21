'use strict';

var app = angular.module('ipl');

/**
 * Fixtures Controller
 * 
 * Controller for fixtures page.
 */
app.controller('fixturesController', ['$http', '$window', '$q', '$mdDialog', 'utilsService', 'urlService', '$scope', '$timeout', function ($http, $window, $q, $mdDialog, utilsService, urlService, $scope, $timeout) {
    var vm = this;
    var token;
    var iNumber;
    vm.init = init;
    vm.searchItem = '';
    vm.makePreditction = makePreditction;
    vm.clearSearchItem = clearSearchItem;
    vm.playerInTeam = playerInTeam;
    vm.checkMOMSelection = checkMOMSelection;
    vm.showMatchStats = showMatchStats;

    vm.updatedTeamVote = updatedTeamVote;
    vm.updatedPlayerVote = updatedPlayerVote;
    vm.updatedCoin = updatedCoin;

    vm.DisbleCard = DisbleCard;

    function DisbleCard(lockPred) {
        if (lockPred) {
            return "lock-pred";
        }
    }

    // Finds if player is part of the playing teams
    function playerInTeam(playerTeamId, teamId1, teamId2) {
        // Akshil -> change this to === later
        if (playerTeamId == teamId1 || playerTeamId == teamId2) {
            return true;
        } else {
            return false;
        }
    }

    function checkMOMSelection(player, predictions) {
        if (predictions)
            return player.id === predictions.momVote;
    }

    // Clears the search bar for select
    function clearSearchItem() {
        vm.searchItem = '';
    }

    // Function to get team object from team id
    function getTeamFromId(teamId) {
        // return vm.teamsList.find(function (team) {
        //     return team.id === teamId;
        // });
        return vm.teamMap[teamId];
    }

    // Init function for main fixtures view
    function init() {
        vm.isLoaded = false;
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
            cache: true,
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        var playersParams = {
            url: urlService.players,
            method: 'GET',
            cache: true,
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        vm.fixturesList = [];
        vm.teamsList = [];
        vm.playersList = [];
        vm.playerMap = {};
        vm.teamMap = {};

        var teamPromise = $http(teamParams);
        var playerPromise = $http(playersParams);
        // Resolve both promises
        $q.all([teamPromise, playerPromise])
            .then(function (data) {
                vm.teamList = data[0].data.teams;
                vm.teamList.forEach(function (team) {
                    vm.teamMap[team.id] = team;
                });
                var role;

                vm.playersList = data[1].data.players;
                vm.playersList.forEach(function (player) {
                    if (player.role === 'allrounder') {
                        player.role = 'All-Rounder';
                    }
                    player.role = utilsService.capitalizeFirstLetter(player.role);

                    vm.playerMap[player.id] = player;
                });

                $http(fixturesParams)
                    .then(function (res) {
                        vm.fixturesList = res.data.matches;
                        vm.fixturesList.forEach(function (fixture) {
                            vm.isLoaded = true;
                            fixture.team1 = getTeamFromId(fixture.teamId1);
                            fixture.team2 = getTeamFromId(fixture.teamId2);
                            fixture.timestamp = moment(fixture.date).format('LLLL');

                            // vm.fixturesList.push({
                            //     teamId1: fixture.teamId1,
                            //     teamId2: fixture.teamId2,
                            //     venue: fixture.venue,
                            //     date: fixture.date,
                            //     timestamp: moment(fixture.date).format('LLLL'),
                            //     status: fixture.status,
                            //     matchId: fixture.id,
                            //     result: fixture.winner,
                            //     manOfMatch: fixture.mom,
                            //     star: fixture.star,
                            //     lockPred: fixture.lock,
                            //     team1: getTeamFromId(fixture.teamId1),
                            //     team2: getTeamFromId(fixture.teamId2),
                            //     predictions: fixture.predictions
                            // });
                        });
                        vm.fixturesList.sort(function (a, b) {
                            return a.id - b.id
                        })
                        vm.fixtureGroup = {};
                        for (var fixture of vm.fixturesList) {
                            var dt = moment(fixture.date).format('dddd LL');
                            if (!vm.fixtureGroup[dt]) {
                                vm.fixtureGroup[dt] = [];
                            }
                            vm.fixtureGroup[dt].push(fixture);
                        }
                    }, function (err) {
                        vm.isLoaded = true;
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
    function updatedTeamVote(teamId, fixture) {
        var data = {
            teamVote: parseInt(teamId),
            inumber: iNumber,
            mid: fixture.matchId,
        }

        if (fixture.predictions && fixture.predictions.predId) {
            if (teamId === fixture.predictions.teamVote) {
                return;
            }
            data.predId = fixture.predictions.predId;
        }
        vm.makePreditction(data, fixture);
    }

    function updatedPlayerVote(playerId, fixture) {
        var data = {
            momVote: playerId,
            inumber: iNumber,
            mid: fixture.matchId,
        }

        if (fixture.predictions && fixture.predictions.predId) {
            if (playerId === fixture.predictions.momVote) {
                return;
            }
            data.predId = fixture.predictions.predId;
        }
        vm.makePreditction(data, fixture);
    }

    function updatedCoin(isUsed, fixture) {
        var data = {
            coinUsed: isUsed,
            inumber: iNumber,
            mid: fixture.matchId,
        }

        if (fixture.predictions && fixture.predictions.predId) {
            if (!isUsed) {
                if (fixture.predictions.coinUsed) {
                    return;
                }
                if (fixture.predictions.coinUsed === false) {
                    return;
                }
            }
            if (isUsed) {
                if (fixture.predictions.coinUsed && fixture.predictions.coinUsed === true) {
                    return;
                }
            }
            data.predId = fixture.predictions.predId;
        }

        vm.makePreditction(data, fixture);
    }



    function makePreditction(data, fixture) {
        var method, url;
        if (!fixture.predictions) {
            method = "POST";
            url = urlService.predictions
        } else {
            method = "PUT";
            url = urlService.predictions + "/" + fixture.predictions.predId;
        }
        var params = {
            url: url,
            method: method,
            data: data,
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        $http(params)
            .then(function (res) {
                if (!fixture.predictions) {
                    fixture.predictions = {};
                    fixture.predictions.predId = res.data.id;
                }
                if (data.teamVote) {
                    fixture.predictions.teamVote = data.teamVote;
                }
                if (data.momVote) {
                    fixture.predictions.momVote = data.momVote;
                }
                if (data.coinUsed) {
                    fixture.predictions.coinUsed = data.coinUsed;
                }
                utilsService.showToast({
                    text: 'Prediction submitted successfully',
                    hideDelay: 0,
                    isError: false
                });
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error: ' + err.data.message, //'Error in submitting prediction, try again later',
                    hideDelay: 0,
                    isError: true
                });
            });
    }

    function showMatchStats(event, id) {
        $mdDialog.show({
            templateUrl: '/static/app/components/fixtures/matchStats.html',
            controller: 'matchStats',
            controllerAs: 'mst',
            targetEvent: event,
            locals: {
                matchId: id,
                teamList: vm.teamsList,
                playerList: vm.playersList,
                playerMap: vm.playerMap,
                teamMap: vm.teamMap
            },
            clickOutsideToClose: true
        }).then(function () {}, function () {});
    }

}]);
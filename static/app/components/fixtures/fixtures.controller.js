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
    vm.flag = {};
    vm.coinIcon = [];
    vm.coinFlag = [];
    vm.selectedTeam = {};
    vm.searchItem = '';
    vm.toggleFlag = toggleFlag;
    vm.useCoin = useCoin;
    vm.makePreditction = makePreditction;
    vm.clearSearchItem = clearSearchItem;
    vm.playerInTeam = playerInTeam;
    vm.openPredictionDialog = openPredictionDialog;
    vm.checkMOMSelection = checkMOMSelection;
    vm.showTeamVote = showTeamVote;
    vm.decideClassName = decideClassName;
    vm.showMatchStats = showMatchStats;

    function decideClassName(predictions) {
        if (predictions && predictions.teamVote) {
            return "voted";
        }
        return "not-voted";
    }

    function showTeamVote(predictions) {
        if (predictions && predictions.teamVote) {
            var teamObj = vm.teamsList.find(function (team) {
                return team.id === predictions.teamVote;
            });
            return "You voted for : " + teamObj.name;
        }
        return "You have not voted for any team yet";

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
            return player.playerId === predictions.momVote;
    }

    // Clears the search bar for select
    function clearSearchItem() {
        vm.searchItem = '';
    }

    // Toggle a flag and select the team chosen to predict
    function toggleFlag(id, teamId1, teamId2,elm) {
        $('.md-fab-fixture,.md-primary').removeClass('highlightTeam');
        $(elm.currentTarget).addClass('highlightTeam');
        vm.flag.id = {};
        vm.flag.id.teamId1 = true;
        vm.flag.id.teamId2 = false;
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
        var statParams = {
            url: urlService.userStats,
            method: 'GET',
            headers: {
                'Accept': 'application/json',
                'Authorization': token
            } 
        }
        vm.fixturesList = [];
        vm.teamsList = [];
        vm.playersList = [];
        vm.predMap = {};
        var playerMap = {};
        var teamMap = {};

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
                    teamMap[team.id]=team;
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
                    playerMap[player.id]=player;
                });

                $http(statParams)
                .then(function(res){
                    res.data.predictions.forEach(function(pred){
                        pred.momN=pred.momVote?playerMap[pred.momVote].name:"-";
                        pred.teamN=pred.teamVote?teamMap[pred.teamVote].shortName:"-";
                        if (!vm.predMap[pred.mid]){
                            vm.predMap[pred.mid]=[];
                        }
                        vm.predMap[pred.mid].push(pred);
                    });
                },function(err){
                    if (err.data.code === 403 && err.data.message === 'token not valid') {
                        utilsService.logout('Session expired, please re-login', true);
                        return;
                    }
                    console.log(err)
                });

                console.log(vm.predMap);



                $http(fixturesParams)
                    .then(function (res) {
                        res.data.matches.forEach(function (fixture) {
                            vm.isLoaded = true;
                            vm.fixturesList.push({
                                teamId1: fixture.teamId1,
                                teamId2: fixture.teamId2,
                                venue: fixture.venue,
                                timestamp: moment(fixture.date).format('LLLL'),
                                status: fixture.status,
                                matchId: fixture.id,
                                result: fixture.winner,
                                manOfMatch: fixture.mom,
                                star: fixture.star,
                                lockPred: fixture.lock,
                                team1: getTeamFromId(fixture.teamId1),
                                team2: getTeamFromId(fixture.teamId2),
                                predictions: fixture.predictions
                            });
                            vm.flag[fixture.mid] = {};
                            vm.coinIcon[fixture.id] = 'money_off';
                            vm.coinFlag[fixture.id] = false;
                        });
                        vm.fixturesList.sort(function (a, b) {
                            return a.matchId - b.matchId
                        })

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
    function makePreditction(fixture) {
        var matchId = fixture.matchId;
        var teamVote = typeof (vm.selectedTeam[matchId]) === "undefined" ? null : vm.selectedTeam[matchId];
        var momVote = typeof (vm.mom[matchId]) === "undefined" ? null : vm.mom[matchId];
        var data = {
            mid: matchId,
            inumber: iNumber,
            teamVote: teamVote,
            momVote: momVote,
            coinUsed: vm.coinFlag[matchId]
        };
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

    function openPredictionDialog(event, fixture) {
        $mdDialog.show({
            templateUrl: '/static/app/components/fixtures/fixturesDialog.html',
            controller: fixturesDialogController,
            scope: $scope,
            preserveScope: true,
            targetEvent: event,
            locals: {
                fixture: fixture
            },
            clickOutesideToClose: true
        });
    }

    function fixturesDialogController($scope, fixture) {
        //var fdc = this;
        // fdc.playersList = vm.playersList;
        $scope.fixture = fixture;

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
                userStats:vm.predMap[id]
            },
            clickOutesideToClose: true
        }).then(function (answer) {
            vm.status = 'You said the information was "' + answer + '".';
        }, function () {
            vm.status = 'You cancelled the dialog.';
        });
    }

}]);
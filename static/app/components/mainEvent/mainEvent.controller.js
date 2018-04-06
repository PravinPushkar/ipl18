'use strict';

var app = angular.module('ipl');

/**
 * Voting Controller
 * 
 * Controller for bonus page.
 */
app.controller('mainEventController', ['$http', '$window', 'urlService', 'utilsService', function ($http, $window, urlService, utilsService) {
    var vm = this;
    var token;
    var teamsQuestionIds = [];

    vm.init = init;
    vm.save = save;
    vm.transformChip = transformChip;
    vm.querySearch = querySearch;
    vm.clearSearchItem = clearSearchItem;
    vm.selectedAnswer = [];

    // Clears the search bar for select box
    function clearSearchItem() {
        vm.searchItem = '';
    }

    // Return the proper object when the append is called.
    function transformChip(chip) {
        if (angular.isObject(chip)) {
            return chip;
        }
        return {
            name: chip,
            type: 'new'
        };
    }

    // Search for the team in autocomplete
    function querySearch(query) {
        var results = query ? vm.teamsListChips.filter(createFilterFor(query)) : [];
        return results;
    }

    // filter function for query search
    function createFilterFor(query) {
        var lowercaseQuery = angular.lowercase(query);

        return function filterFn(team) {
            return (team._lowername.indexOf(lowercaseQuery) === 0);
        };

    }

    // Add _lowercase to teams
    function loadTeamsListChips() {
        return vm.teamsList.map(function (team) {
            team._lowername = team.name.toLowerCase();
            return team;
        });
    }

    // Init function for the main event page
    function init() {
        token = $window.localStorage.getItem('token');
        var questionsParams = {
            url: urlService.bonus,
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

        vm.questions = [];
        vm.teamsList = [];
        vm.playersList = [];
        teamsQuestionIds = [];
        $http(questionsParams)
            .then(function (res) {
                res.data.questions.forEach(function (question) {
                    vm.questions.push({
                        id: question.qid,
                        question: question.question,
                        relatedEntity: question.relatedEntity,
                        points : question.points
                    });
                    if (question.relatedEntity === 'teams') {
                        vm.selectedAnswer[question.qid] = [];
                        teamsQuestionIds.push(question.qid);
                    }
                });
                $http(teamParams)
                    .then(function (res) {
                        res.data.teams.forEach(function (team) {
                            vm.teamsList.push({
                                id: team.id,
                                name: team.name,
                                alias: team.shortName,
                                teamPic: team.picLocation
                            });
                        });
                        vm.teamsListChips = loadTeamsListChips();
                    }, function () {
                        utilsService.showToast({
                            text: 'Error in fetching teams',
                            hideDelay: 0,
                            isError: true
                        });
                    });
                $http(playersParams)
                    .then(function (res) {
                        var role;
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
                    }, function () {
                        utilsService.showToast({
                            text: 'Error in fetching player details',
                            hideDelay: 0,
                            isError: true
                        });
                    });
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
                utilsService.showToast({
                    text: 'Error in fetching main event predictions quiz',
                    hideDelay: 0,
                    isError: true
                });
            });

    }

    // Function will send the answers to backend and save it
    function save(isFormValid) {
        if (isFormValid === false) {
            utilsService.showToast({
                text: 'Please enter valid credentials.',
                hideDelay: 2000,
                isError: true
            });
            return;
        }
        teamsQuestionIds.forEach(function (id) {
            if (vm.selectedAnswer[id].length === 0) {
                utilsService.showToast({
                    text: 'Please enter valid credentials.',
                    hideDelay: 2000,
                    isError: true
                });
                return;
            }
        });

        var dialogParams = {
            title: 'WARNING!',
            text: 'You can submit this only once. You cannot change your answers later.',
            aria: 'Submit Answers',
            ok: 'Continue',
            cancel: 'Cancel'
        };
        utilsService.showConfirmDialog(dialogParams)
            .then(function () {
                var iNumber = $window.localStorage.getItem('iNumber');
                var data = {
                    predictions: []
                };
                vm.questions.forEach(function (question, i) {
                    data.predictions.push({
                        inumber: iNumber,
                        qid: question.id,
                        answer: vm.selectedAnswer[question.id]
                    });
                    if (question.relatedEntity === 'teams') {
                        var tempList = [];
                        vm.selectedAnswer[question.id].forEach(function (item) {
                            tempList.push(item.name);
                        });
                        data.predictions[i].answer = tempList.join(', ');
                    }
                });
                var submitAnswersparams = {
                    url: urlService.bonus,
                    method: 'POST',
                    data: data,
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': token
                    }
                };
                $http(submitAnswersparams)
                    .then(function () {
                        utilsService.showToast({
                            text: 'Answers saved successfully',
                            hideDelay: 1500,
                            isError: false
                        });
                    }, function (err) {
                        if (err.data.code === 403 && err.data.message === 'token not valid') {
                            utilsService.logout('Session expired, please re-login', true);
                            return;
                        }
                        utilsService.showToast({
                            text: 'Error in submitting, please try again later',
                            hideDelay: 0,
                            isError: true
                        });
                    });
            }, function () {
                console.log('Submit cancelled');
            });
    }

}]);
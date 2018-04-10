'use strict';

var app = angular.module('ipl');

/**
 * Fixtures Controller
 * 
 * Controller for fixtures page.
 */
app.controller('matchStats', ['$http','$window','$mdDialog','urlService', 'matchId','teamList','playerList', function ($http,$window,$mdDialog,urlService, matchId,teamList,playerList) {
    var stats = this;
    var token = $window.localStorage.getItem('token');

    stats.init = init;

    stats.hide = function () {
        $mdDialog.hide();
    };

    stats.cancel = function () {
        $mdDialog.cancel();
    };

    stats.close = function () {
        $mdDialog.hide();
    };

    function init() {
        var params = {
            url: `${urlService.fixtures}/${matchId}/stats`,
            method: 'GET',

            headers: {
                'Accept': 'application/json',
                'Authorization': token
            }
        };
        $http(params)
            .then(function (res) {
                stats.teamStats = res.data.teamStats;
                stats.playerStats = res.data.playerStats;
                stats.teamStats = stats.teamStats.reduce(function(acc,curVal,index,arr){
                    curVal.teamName = teamList.filter(function(data){
                        return data.id === curVal.teamId;
                    })[0].name;
                    return arr;
                },[]);
                stats.playerStats = stats.playerStats.reduce(function(acc,curVal,index,arr){
                    curVal.playerName = playerList.filter(function(data){
                        return data.playerId === curVal.playerId;
                    })[0].name;
                    return arr;
                },[]);
            }, function (err) {
                if (err.data.code === 403 && err.data.message === 'token not valid') {
                    utilsService.logout('Session expired, please re-login', true);
                    return;
                }
            });
    }

}]);
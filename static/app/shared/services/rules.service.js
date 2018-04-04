'use strict';

var app = angular.module('ipl');

/**
 * Rules Service
 * 
 * Contains the rules and its subsequent points
 */
app.factory('rulesService', [function () {
    var service = [];

    service.push({
        rule: 'Correct match winner prediction',
        points: '2'
    });
    service.push({
        rule: 'Correct Man of the Match (MoM) prediction',
        points: '4'
    });
    service.push({
        rule: 'For abandoned match',
        points: '1'
    });
    service.push({
        rule: 'Lucky coin used and match is abondoned',
        points: '1 * 5 = 5'
    });
    service.push({
        rule: 'Lucky coin used and correct match winner prediction',
        points: '2 * 5 = 10'
    });
    service.push({
        rule: 'Correct TOP 4 team prediction',
        points: '20'
    });
    service.push({
        rule: 'Correct IPL winner prediction',
        points: '30'
    });
    service.push({
        rule: 'Correct all team standing prediction',
        points: '50'
    });

    return service;
}]);
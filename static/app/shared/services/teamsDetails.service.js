'use strict';

var app = angular.module('ipl');

/**
 * Teams Details
 * 
 * Service that keeps details of a team
 */
app.factory('teamsDetailsService', function () {
    var service = {};

    service.CSK = {
        'venue': 'M. A. Chidambaram Stadium',
        'championships': '2010, 2011',
        'captain': 'MS Dhoni'

    };
    service.DD = {
        'venue': 'Feroz Shah Kotla Grounds',
        'championships': 'None yet',
        'captain': 'Gautam Gambhir'
    };
    service.RR = {
        'venue': 'Sawai Mansingh Stadium',
        'championships': '2008',
        'captain': 'Ajinkya Rahane'
    };
    service.RCB = {
        'venue': 'M. Chinnaswamy Stadium',
        'championships': 'None yet',
        'captain': 'Virat Kohli'
    };
    service.SH = {
        'venue': 'Rajiv Gandhi Intl. Cricket Stadium',
        'championships': '2016',
        'captain': 'Kane Williamson'
    };
    service.KXIP = {
        'venue': 'IS Bindra Stadium',
        'championships': 'None yet',
        'captain': 'Ravichandran Ashwin'
    };
    service.KKR = {
        'venue': 'Eden Gardens',
        'championships': '2012, 2014',
        'captain': 'Dinesh Karthik'
    };
    service.MI = {
        'venue': 'Wankhede Stadium',
        'championships': '2013, 2015, 2017',
        'captain': 'Rohit Sharma'
    };

    return service;
});
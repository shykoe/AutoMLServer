import { combineReducers } from 'redux';
import * as fromPort from './PortReducer';
import * as fromExperiment from './ExperimentReducer'
export const reducer = combineReducers ({PortReducer :  fromPort.PortReducer, ExperimentReducer: fromExperiment.ExperimentIDReducer});
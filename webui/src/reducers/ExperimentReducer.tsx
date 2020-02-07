import { ExperimentIDAction, ExperimentIDTypes } from "../actions/ExperimentAction";
export function ExperimentIDReducer(previousState = `` , action: ExperimentIDAction) : String{
    switch (action.type) {
        case ExperimentIDTypes.Set_ID: 
            const ID = action.payload.ID;
            return ID;
        default:
            return previousState;
    } 
}
import { PortAction, ActionTypes } from "../actions/PortAction";
export function PortReducer(previousState = `` , action: PortAction) : String{
    switch (action.type) {
        case ActionTypes.Set_Port: 
            const port = action.payload.port;
            return port;
        default:
            return previousState;
    } 
}
export enum ActionTypes{
    Set_Port = "webui Set_Port"
}
export interface PortAction { type: ActionTypes.Set_Port, payload: {port: String} };
export function SetPort(port : String): PortAction {
    return {
        type: ActionTypes.Set_Port,
        payload:{
            port:port
        }
    }
};
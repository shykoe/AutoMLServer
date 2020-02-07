export enum ExperimentIDTypes{
    Set_ID = "webui Set_Experiment_ID"
}
export interface ExperimentIDAction { type: ExperimentIDTypes.Set_ID, payload: {ID: String} };
export function SetExperimentID(ID : String): ExperimentIDAction {
    return {
        type: ExperimentIDTypes.Set_ID,
        payload:{
            ID:ID
        }
    }
};
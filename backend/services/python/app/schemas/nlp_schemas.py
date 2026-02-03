from pydantic import BaseModel 



class ComputeCooccurrenceGraph(BaseModel):
    text : str
    window_size : int

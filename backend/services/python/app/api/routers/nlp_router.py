from fastapi import APIRouter
from schemas import ComputeCooccurrenceGraph
from . import nlp_functions as nlp_f
import networkx as nx


router = APIRouter()


@router.post("/compute-graph")
async def compute_cooccurrence_graph(input:ComputeCooccurrenceGraph):

    lang = None
    lang = nlp_f.detect_lang(input.text)

    tokens = nlp_f.preprocess_text_for_nlp_(input.text) if lang == 'ENGLISH' else nlp_f.preprocess_text_for_nlp_fr(input.text)

    G = nx.Graph()

    return { "graph_data" : nlp_f.create_PMI_co_occurrence_graph_for_source(G ,input.window_size,tokens ) }






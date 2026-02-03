import networkx as nx
from itertools import combinations
from networkx.algorithms.community import greedy_modularity_communities
from networkx.algorithms.community.quality import modularity
from collections import defaultdict
from sklearn.preprocessing import minmax_scale
from fastapi import HTTPException
import math
from lingua import Language, LanguageDetectorBuilder
import spacy
import re


nlp = spacy.load("en_core_web_lg")
nlp_fr = spacy.load("fr_core_news_lg")




def detect_lang(raw_text:str,max_chars:int=5000)->str :
    print(" le text :\n",raw_text)
    sample_text = raw_text[:max_chars]
    
    clean_text = ' '.join(sample_text.split())
    
    if len(clean_text) < 50: 
        return None

    languages = [
        Language.ENGLISH, Language.FRENCH, Language.SPANISH, 
        Language.GERMAN, Language.ITALIAN, Language.PORTUGUESE,
    ]
    
    detector = (
        LanguageDetectorBuilder.from_languages(*languages)
        .with_preloaded_language_models()
        .build()
    )

    try:
        language_obj = detector.detect_language_of(clean_text)
        
        if language_obj:
            return language_obj.name 
        else:
            return None

    except Exception as e:
        print(f"compute_reading_path error : {e}")
        import traceback
        traceback.print_exc()  
        return str(e)



def preprocess_text_for_nlp_fr(text: str) -> list[tuple[str, str]]:
    doc = nlp_fr(text)
    results: list[tuple[str, str]] = []
    covered_lemmas: set[str] = set()
    
    
    for chunk in doc.noun_chunks:
        if (
            chunk.root.pos_ in {"NOUN", "PROPN"}
            and chunk.root.is_alpha
            and not chunk.root.is_stop
            and len(chunk) <= 6
        ):
            clean_tokens = [
                token
                for token in chunk
                if token.is_alpha
                and token.pos_  in {"NOUN", "PROPN", "ADJ", "ADP"}
            ]

            if not clean_tokens:
                continue

            clean_chunk = " ".join(
                token.text.lower() for token in clean_tokens
            )

            results.append((clean_chunk, "NOUN"))

            for token in clean_tokens:
                covered_lemmas.add(token.lemma_.lower())
    
    for token in doc:
        lemma = token.lemma_.lower()
        if lemma in covered_lemmas:  
            continue
        if (
            token.pos_ in {"NOUN"}  
            and token.is_alpha
            and not token.is_stop
        ):
            results.append((lemma, token.pos_))
    
    results = [(re.sub(r"[^\w' ]", '', token, flags=re.UNICODE).strip(), _pos) 
               for token, _pos in results]
    
    results = [(token, pos) for token, pos in results if token]
    
    return results



def preprocess_text_for_nlp_(text: str) -> list[tuple[str, str]]:
    doc = nlp(text)
    results: list[tuple[str, str]] = []

    covered_lemmas: set[str] = set()

    for chunk in doc.noun_chunks:
        if (
            chunk.root.pos_ in {"NOUN", "PROPN"}
            and chunk.root.is_alpha
            and not chunk.root.is_stop
            and len(chunk) <= 5
        ):
            clean_tokens = [
                token
                for token in chunk
                if token.is_alpha
                and not token.is_stop
                and token.pos_ != "DET"
                and token.pos_ != "PRON"
            ]

            if not clean_tokens:
                continue

            clean_chunk = " ".join(
                token.text.lower() for token in clean_tokens
            )

            results.append((clean_chunk, "NOUN"))

            for token in clean_tokens:
                covered_lemmas.add(token.lemma_.lower())

    print(f"Après noun chunks: {len(results)} éléments")  # DEBUG
    noun_chunks_count = len(results)


    for token in doc:
        lemma = token.lemma_.lower()

        if lemma in covered_lemmas:
            continue

        if (
            token.pos_ in {"NOUN"}
            and token.is_alpha
            and not token.is_stop
        ):
            results.append((lemma, token.pos_))


    print(f"Après tokens individuels: {len(results)} éléments")  # DEBUG
    
    results = [ (re.sub(r"[^\w' ]", '', token, flags=re.UNICODE).strip(), _pos) for token,_pos in results ]
    
    print(f"Après nettoyage: {len(results)} éléments")  # DEBUG

    print(f"Noun chunks perdus: {noun_chunks_count - len([r for r in results if r[1] == 'NOUN'])}")  # DEBUG

    return results






def h_index(central_node=None , neighbor_strengths=[]):
    if not central_node :
        return 0.0

    if len(list(neighbor_strengths)) == 0 :
           return 0.0

    h = 0
    for i, s in enumerate(neighbor_strengths, start=1):
        if s >= i:
            h = i
        else:
            break

    return h



def compute_neighbors_participation_strengths(u , G  ) :
    try:

        neighbors = list(G.neighbors(u))
        n_neighbors = len(neighbors)
        if n_neighbors == 0:
            return []

        neighbor_strengths = []
        for v in G.neighbors(u):
            w_uv = G[u][v].get("weight",1)
            weights = [G[v][w].get("weight", 1) for w in G.neighbors(v)]
            k = sum(weights)
            if k == 0:
                continue

            v_community_participation_weights = {}
            for w in G.neighbors(v):
                c = G.nodes[w]["modularity_community"] 
                v_w = G[v][w].get("weight", 1)
                v_community_participation_weights[c] = v_community_participation_weights.get(c, 1) + v_w

            v_participation = 1 - sum((weight / k) ** 2 for weight in v_community_participation_weights.values())
            neighbor_strengths.append(v_participation*w_uv)

        if not neighbor_strengths:
           return []

        scaled = minmax_scale(neighbor_strengths) * 10
        return scaled

    except ValueError as e:
        print(f"[compute_participation error] ValueError: {e}")
        raise HTTPException(status_code=500, detail=str(e))

    except TypeError as e:
        print(f"[compute_participation error] TypeError: {e}")
        raise HTTPException(status_code=500, detail=str(e))

    except Exception as e:
        print(f"[compute_participation error] {type(e).__name__}")
        raise HTTPException(status_code=500, detail=str(e))





def compute_node_participation(u,G )->float:

    neighbors = list(G.neighbors(u))
    n_neighbors = len(neighbors)
    if n_neighbors == 0:
        return 0.0


    weights = [G[u][v].get("weight", 1) for v in G.neighbors(u)]
    k = sum(weights)
    
    u_community_participation_weights = {}
    for v in G.neighbors(u):
        c = G.nodes[v]["modularity_community"] 
        u_v = G[u][v].get("weight", 1)
        u_community_participation_weights[c] = u_community_participation_weights.get(c, 0) + u_v

    u_participation = 1 - sum((weight / k) ** 2 for weight in u_community_participation_weights.values())


    return u_participation




def create_PMI_co_occurrence_graph_for_source(G , window_size : int , tokens:list[tuple[str,str]] ) :
    try :

        word_freq = defaultdict(int)
        cooc_freq = defaultdict(int)

        toks = [ t for t , _pos in tokens ] 
        N = len(toks)

        for i in range(len(toks)) :
            word_freq[toks[i]] += 1 

        for i in range(len(toks) - window_size + 1):
            window = toks[i:i + window_size]
            for w1, w2 in combinations(window, 2):
                cooc_freq[tuple(sorted([w1, w2]))] += 1
                
        total_cooc = sum(cooc_freq.values())

        for (w1,w2), cij in cooc_freq.items() :
            
            pi = word_freq[w1] / N
            pj = word_freq[w2] / N
            pij = cij / total_cooc

            if pi == 0 or pj ==  0:

                pmi = float('-inf')
            else:
                pmi = math.log(pij / (pi * pj), 2)

            pmi = max(0, pmi)

                
            if G.has_edge(w1, w2):
                G[w1][w2]["weight"] = cij
                G[w1][w2]["pmi"] = pmi
            else:
                G.add_edge(w1, w2, weight=cij, pmi=pmi)


        edges_to_remove = []
        for u, v, data in G.edges(data=True):
            if data['pmi'] <  4.3:
                edges_to_remove.append((u, v))

        G.remove_edges_from(edges_to_remove)

        isolated_nodes = list(nx.isolates(G))

        G.remove_nodes_from(isolated_nodes)

        G.remove_edges_from(nx.selfloop_edges(G))
        
        betweenness_centrality = nx.betweenness_centrality(G, weight=lambda u, v, d: 1/d['weight'] if d['weight'] > 0 else 1e6)
        clustering = nx.clustering(G , weight="weight")
        page_rank_scores = nx.pagerank(G)

        communities = greedy_modularity_communities(G)

        modularity_indice = modularity(G ,communities )


        nodes_to_remove = []
        for n  in G.nodes():
            if page_rank_scores[n] < 0.0009:
                nodes_to_remove.append(n)


        G.remove_nodes_from(nodes_to_remove)

        node_to_community = {}
        for i , c in enumerate(communities) :
            for node in c :
                node_to_community[node] = i


        node_to_pos = dict(tokens)

        for n in G.nodes:
            G.nodes[n]["_pos"] = node_to_pos.get(n,'')

        for n in G.nodes:
            G.nodes[n]["modularity_community"] = node_to_community.get(n, -1)
        
        for n in G.nodes:
            G.nodes[n]["betweenness"] = betweenness_centrality[n]

        for n in G.nodes:
            G.nodes[n]["page_rank"] = page_rank_scores[n]

        for n in G.nodes:
            G.nodes[n]["freq"] = word_freq[n]


        participation_diversity = {u : h_index(u,compute_neighbors_participation_strengths(u,G)) for u in G.nodes() }
        
        nodes = list(participation_diversity.keys())

        values = [participation_diversity[u] for u in nodes]

        normalized_values = minmax_scale(values, feature_range=(0, 1))

        participation_diversity = dict(zip(nodes, normalized_values))
        
        graph_dict = {
                    "modularity":modularity_indice,
                    "nodes": [{"id": n,"freq": G.nodes[n]["freq"],"pos": G.nodes[n]["_pos"],  "betweenness": betweenness_centrality[n],"page_rank": page_rank_scores[n],"modularity_community": node_to_community[n] , "gatewayness" :  4*(1-clustering[n])+1*betweenness_centrality[n]+ 2 *compute_node_participation(n,G) + 3 * (1 / G.nodes[n]["freq"])   } for n in G.nodes()],
                "edges": [{"source": u, "target": v, "weight": d["weight"] , "pmi" : d["pmi"] }  for u, v, d in G.edges(data=True)]
            }

       
        return graph_dict


    
    except ValueError as e:
        print(f"[create_graph error] ValueError: {e}")
        import traceback
        traceback.print_exc()  # ← IMPORTANT: Affiche la stack trace complète
        print(str(e))
        raise HTTPException(status_code=500, detail=str(e))


    except TypeError as e:
        import traceback
        traceback.print_exc()  # ← IMPORTANT: Affiche la stack trace complète
        print(str(e))
        raise HTTPException(status_code=500, detail=str(e))


    except Exception as e:
        import traceback
        traceback.print_exc()  # ← IMPORTANT: Affiche la stack trace complète
        print(str(e))
        raise HTTPException(status_code=500, detail=str(e))


    except Exception as e :
        import traceback
        traceback.print_exc()  # ← IMPORTANT: Affiche la stack trace complète
        print(str(e))
        raise HTTPException(status_code=500, detail=str(e))






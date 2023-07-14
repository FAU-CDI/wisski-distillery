# This file is used to initialize a new GraphDB repository. 
# It is filled at runtime.

@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#>.
@prefix rep: <http://www.openrdf.org/config/repository#>.
@prefix sr: <http://www.openrdf.org/config/repository/sail#>.
@prefix sail: <http://www.openrdf.org/config/sail#>.
@prefix owlim: <http://www.ontotext.com/trree/owlim#>.

[] a rep:Repository ;
    rep:repositoryID "{{ .RepositoryID }}" ;
    rdfs:label "{{ .Label }}" ;
    rep:repositoryImpl [
        rep:repositoryType "graphdb:SailRepository" ;
        sr:sailImpl [
            sail:sailType "graphdb:Sail" ;

            owlim:owlim-license "" ;

            owlim:base-URL "{{ .BaseURL }}" ;
            owlim:defaultNS "" ;
            owlim:entity-index-size "10000000" ;
            owlim:entity-id-size  "32" ;
            owlim:imports "" ;
            owlim:repository-type "file-repository" ;
            owlim:ruleset "empty" ;
            owlim:storage-folder "storage" ;

            owlim:enable-context-index "true" ;
            owlim:cache-memory "80m" ;
            owlim:tuple-index-memory "80m" ;

            owlim:enablePredicateList "true" ;
            owlim:predicate-memory "0%" ;

            owlim:fts-memory "0%" ;
            owlim:ftsIndexPolicy "never" ;
            owlim:ftsLiteralsOnly "true" ;

            owlim:in-memory-literal-properties "false" ;
            owlim:enable-literal-index "true" ;
            owlim:index-compression-ratio "-1" ;

            owlim:check-for-inconsistencies "false" ;
            owlim:disable-sameAs  "false" ;
            owlim:enable-optimization  "true" ;
            owlim:transaction-mode "safe" ;
            owlim:transaction-isolation "true" ;
            owlim:query-timeout  "0" ;
            owlim:query-limit-results  "0" ;
            owlim:throw-QueryEvaluationException-on-timeout "false" ;
            owlim:useShutdownHooks  "true" ;
            owlim:read-only "false" ;
            owlim:nonInterpretablePredicates "http://www.w3.org/2000/01/rdf-schema#label;http://www.w3.org/1999/02/22-rdf-syntax-ns#type;http://www.ontotext.com/owlim/ces#gazetteerConfig;http://www.ontotext.com/owlim/ces#metadataConfig" ;
        ]
    ].
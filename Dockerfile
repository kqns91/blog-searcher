FROM opensearchproject/opensearch:1.1.0
RUN /usr/share/opensearch/bin/opensearch-plugin install analysis-kuromoji 
RUN /usr/share/opensearch/bin/opensearch-plugin install analysis-icu
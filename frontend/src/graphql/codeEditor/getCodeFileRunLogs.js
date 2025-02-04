import { gql, GraphQLClient } from 'graphql-request';
import { useGlobalAuthState } from '../../Auth/UserAuth.jsx';

const graphlqlEndpoint = import.meta.env.VITE_GRAPHQL_ENDPOINT_PRIVATE;

const query = gql`
    query getCodeFileRunLogs($environmentID: String!, $pipelineID: String!, $runID: String!) {
        getCodeFileRunLogs(environmentID: $environmentID, pipelineID: $pipelineID, runID: $runID) {
            created_at
            uid
            log
            log_type
        }
    }
`;

export const useGetCodeFileRunLogs = () => {
    const authState = useGlobalAuthState();
    const jwt = authState.authToken.get();

    const headers = {
        Authorization: 'Bearer ' + jwt,
    };

    const client = new GraphQLClient(graphlqlEndpoint, {
        headers,
    });

    return async (input) => {
        try {
            const res = await client.request(query, input);
            return res?.getCodeFileRunLogs;
        } catch (error) {
            return JSON.parse(JSON.stringify(error, undefined, 2)).response;
        }
    };
};

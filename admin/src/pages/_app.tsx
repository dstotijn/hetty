import { ApolloProvider } from "@apollo/client";
import { CacheProvider, EmotionCache } from "@emotion/react";
import { ThemeProvider } from "@mui/material";
import CssBaseline from "@mui/material/CssBaseline";
import { AppProps } from "next/app";
import Head from "next/head";
import React from "react";

import { ActiveProjectProvider } from "lib/ActiveProjectContext";
import { InterceptedRequestsProvider } from "lib/InterceptedRequestsContext";
import { useApollo } from "lib/graphql/useApollo";
import createEmotionCache from "lib/mui/createEmotionCache";
import theme from "lib/mui/theme";

import "../styles.css";

// Client-side cache, shared for the whole session of the user in the browser.
const clientSideEmotionCache = createEmotionCache();

interface MyAppProps extends AppProps {
  emotionCache?: EmotionCache;
}

export default function MyApp(props: MyAppProps) {
  const { Component, emotionCache = clientSideEmotionCache, pageProps } = props;
  const apolloClient = useApollo();

  return (
    <CacheProvider value={emotionCache}>
      <Head>
        <title>Hetty://</title>
        <meta name="viewport" content="initial-scale=1, width=device-width" />
      </Head>
      <ApolloProvider client={apolloClient}>
        <ActiveProjectProvider>
          <InterceptedRequestsProvider>
            <ThemeProvider theme={theme}>
              <CssBaseline />
              <Component {...pageProps} />
            </ThemeProvider>
          </InterceptedRequestsProvider>
        </ActiveProjectProvider>
      </ApolloProvider>
    </CacheProvider>
  );
}

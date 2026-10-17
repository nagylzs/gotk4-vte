#pragma once

#include <stdlib.h>

static inline char ** make_strings(int count)
{
    return (char **)malloc(sizeof(char *) * count);
}

static inline void destroy_strings(char **strings, int count)
{
    for (int i = 0; i < count; i++) {
        free(strings[i]);
    }
    free(strings);
}

static inline void set_string(char **strings, int n, char *str)
{
    strings[n] = str;
}

static VteTerminal *toVteTerminal(void *p)
{
    return (VTE_TERMINAL(p));
}

static inline gboolean set_context_menu_model(VteTerminal *terminal, GMenuModel *model)
{
#if VTE_CHECK_VERSION(0, 76, 0)
    vte_terminal_set_context_menu_model(terminal, model);
    return TRUE;
#else
    return FALSE;
#endif
}

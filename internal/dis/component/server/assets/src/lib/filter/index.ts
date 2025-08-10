document.querySelectorAll<HTMLInputElement>('input[data-filter-list]').forEach(element => {
    const filterListAttribute = element.dataset.filterList;
    if (
        typeof filterListAttribute !== 'string' ||
        !/^[a-zA-Z_-]+$/.test(filterListAttribute)
    ) {
        console.warn('Element', element, 'has invalid or missing data-filter-list attribute', filterListAttribute);
        return;
    }

    element.form?.addEventListener('submit', (ev) => ev.preventDefault());

    const selector = '*[' + CSS.escape(filterListAttribute) + ']';

    /**
     * Check if the given element matches a query
     * @param element Element to check
     * @param lowercaseQuery Query to match. Assumed to be lowercase.
     * @returns 
     */
    const matches = (element: HTMLElement, lowercaseQuery: string): boolean => {
        if (lowercaseQuery.length === 0) {
            return true;
        }

        const value = element.getAttribute(filterListAttribute)?.toLowerCase();
        if (typeof value !== 'string') {
            return true;
        }
        
        // simple fuzzy match
        let j = 0;
        const len = lowercaseQuery.length;
        for (let i = 0; i < value.length && j < len; i++) {
            if (value[i] !== lowercaseQuery[j]) {
                continue;
            }
            j++;
        }
        return j === len;
    }

    /**
     * Updates the view to display only elements matching the query.
     * @param lowercaseQuery Query to match against. Assumed to be lowercase.
     */
    const updateFilter = (lowercaseQuery: string) => {
        document
            .querySelectorAll<HTMLElement>(selector)
            .forEach(candidate => {
                if (matches(candidate, lowercaseQuery)) {
                    candidate.style.display = '';
                } else {
                    candidate.style.display = 'none';
                }
            });
    };

    element.addEventListener('input', (ev) => {
        ev.preventDefault();
        const target = ev.currentTarget as HTMLInputElement;
        updateFilter(target.value.toLowerCase());
    });
});

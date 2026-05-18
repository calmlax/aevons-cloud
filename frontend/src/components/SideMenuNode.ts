import { defineComponent, h, resolveComponent, type PropType } from 'vue';
import { useI18n } from 'vue-i18n';
import SvgIcon from './SvgIcon/index.vue';

import type { MenuNode } from '../types/menu';

const SideMenuNode = defineComponent({
  name: 'SideMenuNode',
  props: {
    item: {
      type: Object as PropType<MenuNode>,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();

    return () => {
      // Filter out explicitly hidden menu items
      if (props.item.meta?.hidden === true) {
        return null;
      }

      const ASubMenu = resolveComponent('ASubMenu');
      const AMenuItem = resolveComponent('AMenuItem');
      const iconName = props.item.meta?.icon;
      const iconSlot = iconName ? () => h(SvgIcon, { name: iconName }) : undefined;

      const titleKey = props.item.meta?.titleKey || props.item.meta?.title || 'layout.defaultTitle';

      if (props.item.children?.length) {
        return h(
          ASubMenu,
          { key: props.item.key },
          {
            icon: iconSlot,
            title: () => t(titleKey),
            default: () => props.item.children!.map((child) => h(SideMenuNode, { key: child.key, item: child })),
          }
        );
      }

      return h(
        AMenuItem,
        { key: props.item.key },
        {
          icon: iconSlot,
          default: () => t(titleKey),
        }
      );
    };
  },
});

export default SideMenuNode;
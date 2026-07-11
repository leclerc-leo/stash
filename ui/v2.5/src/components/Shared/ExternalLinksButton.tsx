import { Button, Dropdown } from "react-bootstrap";
import { ExternalLink } from "./ExternalLink";
import TextUtils from "src/utils/text";
import { Icon } from "./Icon";
import { IconDefinition, faLink, faBox, faServer , faPaw, faGamepad, faCity, faBoxArchive } from "@fortawesome/free-solid-svg-icons";
import { useMemo } from "react";
import { faInstagram, faTwitter, faPatreon } from "@fortawesome/free-brands-svg-icons";
import ReactDOM from "react-dom";
import { PatchComponent } from "src/patch";

export const ExternalLinksButton: React.FC<{
  icon?: IconDefinition;
  urls: string[];
  className?: string;
  openIfSingle?: boolean;
  onClick?: (url: string) => void;
}> = PatchComponent(
  "ExternalLinksButton",
  ({ urls, icon = faLink, className = "", openIfSingle = false, onClick }) => {
    if (!urls.length) {
      return null;
    }

    const Menu = () =>
      ReactDOM.createPortal(
        <Dropdown.Menu>
          {urls.map((url) => (
            <Dropdown.Item
              key={url}
              as={ExternalLink}
              href={TextUtils.sanitiseURL(url)}
              title={url}
              onClick={(e) => {
                onClick?.(url);
              }}
            >
              {url}
            </Dropdown.Item>
          ))}
        </Dropdown.Menu>,
        document.body
      );

    if (openIfSingle && urls.length === 1) {
      return (
        <ExternalLink
          className={`external-links-button-link minimal btn link ${className}`}
          href={TextUtils.sanitiseURL(urls[0])}
          target={urls[0].startsWith("sftp:") || urls[0].startsWith("ftp:") ? undefined : "_blank"}
          rel="noopener noreferrer"
          onClick={(e) => {
            onClick?.(urls[0]);
          }}
        >
          <Icon icon={icon} />
        </ExternalLink>
      );
    } else {
      return (
        <Dropdown className="external-links-button">
          <Dropdown.Toggle as={Button} className={`minimal link ${className}`}>
            <Icon icon={icon} />
          </Dropdown.Toggle>
          <Menu />
        </Dropdown>
      );
    }
  }
);

export const ExternalLinkButtons: React.FC<{ urls: string[] | undefined, onClick?: (url: string) => void }> =
  PatchComponent("ExternalLinkButtons", ({ urls, onClick }) => {
    const urlSpecs = useMemo(() => {
      if (!urls?.length) {
        return [];
      }

      const twitter = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?(?:twitter|x).com\//)
      );
      const instagram = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?instagram.com\//)
      );
      const patreon = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?patreon.com\//)
      );
      const fanbox = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?[^\/]+\.fanbox\.cc\/posts\//)
      );
      const pawchive = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?pawchive\./)
      );
      const f95zone = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?f95zone\.to\//)
      );
      const simpcity = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?simpcity\./)
      );
      const coomer = urls.filter((u) =>
        u.match(/https?:\/\/(?:www\.)?coomer\./)
      );
      const sftp = urls.filter((u) => u.startsWith("sftp://") || u.startsWith("ftp://"));
      const others = urls.filter(
        (u) => !twitter.includes(u) && !instagram.includes(u) && !patreon.includes(u) && !fanbox.includes(u)
            && !pawchive.includes(u) && !f95zone.includes(u) && !simpcity.includes(u)
            && !coomer.includes(u) && !sftp.includes(u)
      );

      return [
        { icon: faServer, className: "sftp", urls: sftp },
        { icon: faGamepad, className: "f95zone", urls: f95zone },
        { icon: faCity, className: "simpcity", urls: simpcity },
        { icon: faPatreon, className: "patreon", urls: patreon },
        { icon: faBox, className: "fanbox", urls: fanbox },
        { icon: faPaw, className: "pawchive", urls: pawchive },
        { icon: faBoxArchive, className: "coomer", urls: coomer },
        { icon: faLink, className: "", urls: others },
        { icon: faTwitter, className: "twitter", urls: twitter },
        { icon: faInstagram, className: "instagram", urls: instagram },
      ];
    }, [urls]);

    return (
      <>
        {urlSpecs.map((spec, i) => (
          <ExternalLinksButton key={i} {...spec} onClick={onClick} openIfSingle />
        ))}
      </>
    );
  });

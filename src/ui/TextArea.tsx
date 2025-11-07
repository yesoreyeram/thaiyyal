type Props = {
  text: string;
  placeholder?: string;
  rows?: number;
  onChange: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
};

export const TextArea = (props: Props) => {
  const { text, onChange, placeholder = "Enter text...", rows = 3 } = props;
  return (
    <textarea
      value={String(text ?? "")}
      onChange={onChange}
      rows={rows}
      placeholder={placeholder}
      aria-label="Text input value"
      className="w-full text-xs border px-1.5 py-0.5 rounded  focus:ring-1 focus:outline-none border-gray-600  bg-gray-900 text-white focus:ring-blue-400 "
    />
  );
};
